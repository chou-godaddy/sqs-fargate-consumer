package consumer

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

type MetricsCollector struct {
	client       *cloudwatch.Client
	metrics      chan Metric
	metricValues map[string]float64
	metricWindow map[string][]MetricDataPoint
	mu           sync.RWMutex
	namespace    string
}

type Metric struct {
	Name      string
	Value     float64
	Unit      string
	Timestamp time.Time
}

type MetricDataPoint struct {
	Value     float64
	Timestamp time.Time
}

const (
	MetricWorkerCount       = "WorkerCount"
	MetricWorkerUtilization = "WorkerUtilization"
	MetricQueueDepth        = "QueueDepth"
	MetricProcessingTime    = "ProcessingTime"
)

func NewMetricsCollector(client *cloudwatch.Client, namespace string) *MetricsCollector {
	return &MetricsCollector{
		client:       client,
		metrics:      make(chan Metric, 1000),
		metricValues: make(map[string]float64),
		metricWindow: make(map[string][]MetricDataPoint),
		namespace:    namespace,
	}
}

// GetMetric returns the current value of a metric
func (c *MetricsCollector) GetMetric(name string) float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Return the current value if it exists
	if value, exists := c.metricValues[name]; exists {
		return value
	}

	return 0
}

// GetMetricAverage returns the average value over the specified duration
func (c *MetricsCollector) GetMetricAverage(name string, duration time.Duration) float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	dataPoints := c.metricWindow[name]
	if len(dataPoints) == 0 {
		return 0
	}

	cutoff := time.Now().Add(-duration)
	var sum float64
	var count int

	for _, dp := range dataPoints {
		if dp.Timestamp.After(cutoff) {
			sum += dp.Value
			count++
		}
	}

	if count == 0 {
		return 0
	}

	return sum / float64(count)
}

// RecordMetric records a new metric value
func (c *MetricsCollector) RecordMetric(name string, value float64, unit string) {
	metric := Metric{
		Name:      name,
		Value:     value,
		Unit:      unit,
		Timestamp: time.Now(),
	}

	c.metrics <- metric
}

// processMetric updates internal metric state
func (c *MetricsCollector) processMetric(metric Metric) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Update current value
	c.metricValues[metric.Name] = metric.Value

	// Add to sliding window
	dataPoint := MetricDataPoint{
		Value:     metric.Value,
		Timestamp: metric.Timestamp,
	}

	// Initialize slice if it doesn't exist
	if _, exists := c.metricWindow[metric.Name]; !exists {
		c.metricWindow[metric.Name] = make([]MetricDataPoint, 0)
	}

	// Add new datapoint
	c.metricWindow[metric.Name] = append(c.metricWindow[metric.Name], dataPoint)

	// Cleanup old datapoints (keep last hour)
	c.cleanupOldDataPoints(metric.Name)
}

// cleanupOldDataPoints removes datapoints older than 30 minutes
func (c *MetricsCollector) cleanupOldDataPoints(metricName string) {
	cutoff := time.Now().Add(-30 * time.Minute)
	dataPoints := c.metricWindow[metricName]

	// Find index of first datapoint to keep
	keepIndex := 0
	for i, dp := range dataPoints {
		if dp.Timestamp.After(cutoff) {
			keepIndex = i
			break
		}
	}

	// Slice off old datapoints
	if keepIndex > 0 {
		c.metricWindow[metricName] = dataPoints[keepIndex:]
	}
}

// Start begins processing metrics and publishing to CloudWatch
func (c *MetricsCollector) Start(ctx context.Context) {
	// Process incoming metrics
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case metric := <-c.metrics:
				c.processMetric(metric)
			}
		}
	}()

	// Publish to CloudWatch periodically
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.publishMetrics()
		}
	}
}

// publishMetrics publishes all current metrics to CloudWatch
func (c *MetricsCollector) publishMetrics() {
	c.mu.RLock()
	metricData := make([]types.MetricDatum, 0)

	for name, value := range c.metricValues {
		metricData = append(metricData, types.MetricDatum{
			MetricName: &name,
			Value:      &value,
			Timestamp:  aws.Time(time.Now()),
			Unit:       types.StandardUnitCount, // Adjust based on metric type
		})
	}
	c.mu.RUnlock()

	if len(metricData) > 0 {
		_, err := c.client.PutMetricData(context.Background(), &cloudwatch.PutMetricDataInput{
			Namespace:  &c.namespace,
			MetricData: metricData,
		})
		if err != nil {
			log.Printf("Error publishing metrics to CloudWatch: %v", err)
		}
	}
}

// RecordProcessingTime records message processing duration
func (c *MetricsCollector) RecordProcessingTime(duration time.Duration) {
	c.RecordMetric(MetricProcessingTime, float64(duration.Milliseconds()), "Milliseconds")
}

// RecordQueueDepth records current SQS queue depth
func (c *MetricsCollector) RecordQueueDepth(depth int) {
	c.RecordMetric(MetricQueueDepth, float64(depth), "Count")
}

// RecordWorkerUtilization records worker utilization percentage
func (c *MetricsCollector) RecordWorkerUtilization(utilizationPercentage float64) {
	c.RecordMetric(MetricWorkerUtilization, utilizationPercentage, "Percent")
}

// RecordError records an error occurrence
func (c *MetricsCollector) RecordError(errorType string) {
	c.RecordMetric("Error_"+errorType, 1, "Count")
}

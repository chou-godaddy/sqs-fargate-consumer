package consumer

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/google/uuid"
)

type Config struct {
	QueueURL             string
	InitialWorkerCount   int
	MaxWorkerCount       int
	MinWorkerCount       int
	ScaleUpThreshold     float64
	ScaleDownThreshold   float64
	CooldownPeriod       time.Duration
	TargetProcessingTime time.Duration
}

type Consumer struct {
	client           *sqs.Client
	queueURL         string
	workers          map[string]*Worker
	metricsCollector *MetricsCollector
	mu               sync.RWMutex
	config           *Config
}

func NewConsumer(client *sqs.Client, collector *MetricsCollector, config *Config) *Consumer {
	return &Consumer{
		client:           client,
		queueURL:         config.QueueURL,
		workers:          make(map[string]*Worker),
		metricsCollector: collector,
		config:           config,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	// Start initial workers
	for i := 0; i < c.config.InitialWorkerCount; i++ {
		if err := c.addWorker(); err != nil {
			return fmt.Errorf("failed to start initial workers: %w", err)
		}
	}

	// Start metrics reporting
	go c.reportMetrics(ctx)

	// Start queue depth monitoring
	go c.monitorQueueDepth(ctx)

	return nil
}

func (c *Consumer) monitorQueueDepth(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			attrs, err := c.client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
				QueueUrl: &c.queueURL,
				AttributeNames: []types.QueueAttributeName{
					types.QueueAttributeNameApproximateNumberOfMessages,
					types.QueueAttributeNameApproximateNumberOfMessagesNotVisible,
				},
			})

			if err != nil {
				c.metricsCollector.RecordError("queue_depth_fetch_error")
				continue
			}

			// Get visible messages
			if visibleStr, ok := attrs.Attributes[string(types.QueueAttributeNameApproximateNumberOfMessages)]; ok {
				visible, err := strconv.Atoi(visibleStr)
				if err == nil {
					c.metricsCollector.RecordQueueDepth(visible)
				} else {
					c.metricsCollector.RecordError("queue_depth_parse_error")
				}
			}

			// Get in-flight messages
			if notVisibleStr, ok := attrs.Attributes[string(types.QueueAttributeNameApproximateNumberOfMessagesNotVisible)]; ok {
				notVisible, err := strconv.Atoi(notVisibleStr)
				if err == nil {
					c.metricsCollector.RecordMetric("InFlightMessages", float64(notVisible), "Count")
				} else {
					c.metricsCollector.RecordError("in_flight_messages_parse_error")
				}
			}
		}
	}
}

func (c *Consumer) addWorker() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.workers) >= c.config.MaxWorkerCount {
		return fmt.Errorf("max worker count reached")
	}

	workerID := uuid.New().String()
	worker := NewWorker(workerID, c.client, c.queueURL, c.metricsCollector)
	c.workers[workerID] = worker

	go worker.Start(context.Background())

	return nil
}

func (c *Consumer) removeWorker() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.workers) <= c.config.MinWorkerCount {
		return fmt.Errorf("min worker count reached")
	}

	// Remove least active worker
	var workerToRemove string
	for id := range c.workers {
		workerToRemove = id
		break
	}

	if worker, exists := c.workers[workerToRemove]; exists {
		worker.Stop()
		delete(c.workers, workerToRemove)
	}

	return nil
}

func (c *Consumer) reportMetrics(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.mu.RLock()
			workerCount := len(c.workers)
			activeWorkers := 0
			for _, worker := range c.workers {
				if worker.isProcessing.Load() {
					activeWorkers++
				}
			}
			c.mu.RUnlock()

			c.metricsCollector.RecordMetric(MetricWorkerCount, float64(workerCount), "Count")

			if workerCount > 0 {
				utilization := (float64(activeWorkers) / float64(workerCount)) * 100
				c.metricsCollector.RecordWorkerUtilization(utilization)
			}
		}
	}
}

func (c *Consumer) Shutdown() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, worker := range c.workers {
		worker.Stop()
	}
}

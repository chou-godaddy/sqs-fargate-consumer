package consumer

import (
	"context"
	"log"
	"time"
)

type Scaler struct {
	collector      *MetricsCollector
	consumer       *Consumer
	config         *Config
	lastScaleEvent time.Time
}

func NewScaler(collector *MetricsCollector, consumer *Consumer, config *Config) *Scaler {
	return &Scaler{
		collector:      collector,
		consumer:       consumer,
		config:         config,
		lastScaleEvent: time.Now(),
	}
}

func (s *Scaler) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.evaluateScaling()
		}
	}
}

func (s *Scaler) evaluateScaling() {
	queueDepth := s.collector.GetMetric(MetricQueueDepth)
	processingTime := s.collector.GetMetricAverage(MetricProcessingTime, 1*time.Minute)
	workerUtilization := s.collector.GetMetricAverage(MetricWorkerUtilization, 1*time.Minute)

	if s.shouldScaleUp(queueDepth, processingTime, workerUtilization) {
		err := s.consumer.addWorker()
		if err != nil {
			log.Printf("Failed to scale up: %v", err)
		}
		s.lastScaleEvent = time.Now()
		return
	}

	if s.shouldScaleDown(queueDepth, processingTime, workerUtilization) {
		if time.Since(s.lastScaleEvent) > s.config.CooldownPeriod {
			err := s.consumer.removeWorker()
			if err != nil {
				log.Printf("Failed to scale down: %v", err)
			}
			s.lastScaleEvent = time.Now()
		}
	}
}

func (s *Scaler) shouldScaleUp(queueDepth, processingTime, utilization float64) bool {
	return queueDepth > s.config.ScaleUpThreshold ||
		utilization > 75.0 ||
		processingTime > s.config.TargetProcessingTime.Seconds()
}

func (s *Scaler) shouldScaleDown(queueDepth, processingTime, utilization float64) bool {
	return queueDepth < s.config.ScaleDownThreshold &&
		utilization < 40.0 &&
		processingTime < s.config.TargetProcessingTime.Seconds()
}

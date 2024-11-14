package consumer

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Worker struct {
	id               string
	client           *sqs.Client
	queueURL         string
	metricsCollector *MetricsCollector
	done             chan bool
	processingTime   time.Duration
	isProcessing     atomic.Bool  // Track if worker is currently processing
	processedCount   atomic.Int64 // Count of messages processed in current window
	lastWindowTime   atomic.Int64 // Last time window for utilization calculation
}

func NewWorker(id string, client *sqs.Client, queueURL string, collector *MetricsCollector) *Worker {
	w := &Worker{
		id:               id,
		client:           client,
		queueURL:         queueURL,
		metricsCollector: collector,
		done:             make(chan bool),
		processingTime:   30 * time.Second,
	}

	// Initialize atomic values
	w.lastWindowTime.Store(time.Now().Unix())

	// Start utilization tracking
	go w.trackUtilization()

	return w
}

func (w *Worker) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-w.done:
			return
		default:
			messages, err := w.pollMessages(ctx)
			if err != nil {
				w.metricsCollector.RecordError("poll_error")
				time.Sleep(1 * time.Second) // Backoff on error
				continue
			}

			for _, msg := range messages {
				startTime := time.Now()
				w.isProcessing.Store(true)

				if err := w.processMessage(ctx, msg); err != nil {
					w.handleError(ctx, msg, err)
				} else {
					w.deleteMessage(ctx, msg)
					w.processedCount.Add(1)
				}

				w.isProcessing.Store(false)
				processingDuration := time.Since(startTime)
				w.metricsCollector.RecordProcessingTime(processingDuration)
			}
		}
	}
}

func (w *Worker) trackUtilization() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-w.done:
			return
		case <-ticker.C:
			now := time.Now().Unix()
			lastWindow := w.lastWindowTime.Swap(now)

			// Calculate messages processed per second in the window
			timeWindow := now - lastWindow
			if timeWindow > 0 {
				processed := w.processedCount.Swap(0)
				utilization := float64(processed) / float64(timeWindow) * 100

				// Cap utilization at 100%
				if utilization > 100 {
					utilization = 100
				}

				w.metricsCollector.RecordWorkerUtilization(utilization)
			}
		}
	}
}

func (w *Worker) Stop() {
	close(w.done)
}

func (w *Worker) pollMessages(ctx context.Context) ([]types.Message, error) {
	output, err := w.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            &w.queueURL,
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     20, // Long polling
		VisibilityTimeout:   int32(w.processingTime.Seconds()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to poll messages: %w", err)
	}
	return output.Messages, nil
}

func (w *Worker) processMessage(ctx context.Context, msg types.Message) error {
	log.Printf("Processing message %s", *msg.MessageId)
	time.Sleep(5 * time.Second) // Simulate processing time
	return nil
}

func (w *Worker) deleteMessage(ctx context.Context, msg types.Message) error {
	_, err := w.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &w.queueURL,
		ReceiptHandle: msg.ReceiptHandle,
	})
	return err
}

func (w *Worker) handleError(ctx context.Context, msg types.Message, err error) {
	log.Printf("Error processing message %s: %v", *msg.MessageId, err)
	w.metricsCollector.RecordError("processing_error")

	// Modify visibility timeout to retry later
	_, changeErr := w.client.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
		QueueUrl:          &w.queueURL,
		ReceiptHandle:     msg.ReceiptHandle,
		VisibilityTimeout: 30, // Reset to 30 seconds
	})
	if changeErr != nil {
		log.Printf("Error changing message visibility: %v", changeErr)
	}
}

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sqs-fargate-consumer/internal/consumer"
	"sqs-fargate-consumer/internal/utils"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load AWS configuration
	awscfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	cfg := &consumer.Config{
		QueueURL:             "https://sqs.us-west-2.amazonaws.com/982600293865/sqs-fargate-consumer-eventqueue",
		InitialWorkerCount:   5,
		MaxWorkerCount:       100,
		MinWorkerCount:       2,
		ScaleUpThreshold:     50,
		ScaleDownThreshold:   10,
		CooldownPeriod:       30 * time.Second,
		TargetProcessingTime: 30 * time.Second,
	}

	// Health and ready check server
	go utils.StartHTTPSServer()

	// Initialize components
	sqsClient := sqs.NewFromConfig(awscfg)
	cwClient := cloudwatch.NewFromConfig(awscfg)
	metricsCollector := consumer.NewMetricsCollector(cwClient, "SQSConsumer")
	sqsConsumer := consumer.NewConsumer(sqsClient, metricsCollector, cfg)
	scaler := consumer.NewScaler(metricsCollector, sqsConsumer, cfg)

	// Start the consumer system
	go sqsConsumer.Start(ctx)
	go metricsCollector.Start(ctx)
	go scaler.Start(ctx)

	// Wait for shutdown signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	<-sigCh

	// Graceful shutdown
	cancel()
	sqsConsumer.Shutdown()
}

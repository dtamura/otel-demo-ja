// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"log"
	"os"
	"strings"
)

func main() {
	log.Println("Accounting service started")

	// Output relevant environment variables
	envs := filterRelevantEnvVars()
	outputEnvVarsInOrder(envs)

	ctx := context.Background()

	// Setup tracing
	tp, err := setupTracing(ctx)
	if err != nil {
		log.Fatalf("Failed to setup tracing: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Failed to shutdown tracer provider: %v", err)
		}
	}()

	// Setup database
	db, err := setupDatabase()
	if err != nil {
		log.Fatalf("Failed to setup database: %v", err)
	}

	// Setup Kafka consumer
	kafkaAddr := os.Getenv("KAFKA_ADDR")
	if kafkaAddr == "" {
		log.Fatal("KAFKA_ADDR environment variable is required")
	}

	brokers := strings.Split(kafkaAddr, ",")
	log.Printf("Connecting to Kafka: %s", kafkaAddr)

	consumer, err := NewConsumer(brokers, db)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	// Start consuming messages
	if err := consumer.StartListening(ctx); err != nil {
		log.Fatalf("Failed to start consuming: %v", err)
	}
}

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"

	oteldemo "accounting-go/genproto"
)

const topicName = "orders"

// Consumer handles Kafka message consumption
type Consumer struct {
	consumer sarama.Consumer
	db       *gorm.DB
	tracer   trace.Tracer
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(brokers []string, db *gorm.DB) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	tracer := otel.Tracer("Accounting.Consumer")

	return &Consumer{
		consumer: consumer,
		db:       db,
		tracer:   tracer,
	}, nil
}

// StartListening starts consuming messages from Kafka
func (c *Consumer) StartListening(ctx context.Context) error {
	partitionConsumer, err := c.consumer.ConsumePartition(topicName, 0, sarama.OffsetOldest)
	if err != nil {
		return fmt.Errorf("failed to create partition consumer: %w", err)
	}
	defer partitionConsumer.Close()

	log.Printf("Starting to consume from topic: %s", topicName)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case message := <-partitionConsumer.Messages():
			c.processMessage(ctx, message)
		case err := <-partitionConsumer.Errors():
			log.Printf("Consumer error: %v", err)
		case <-signals:
			log.Println("Interrupt signal received, shutting down...")
			return nil
		case <-ctx.Done():
			log.Println("Context cancelled, shutting down...")
			return nil
		}
	}
}

// processMessage processes a Kafka message
func (c *Consumer) processMessage(ctx context.Context, message *sarama.ConsumerMessage) {
	ctx, span := c.tracer.Start(ctx, "order-consumed",
		trace.WithSpanKind(trace.SpanKindInternal),
		trace.WithAttributes(
			attribute.String("kafka.topic", message.Topic),
			attribute.Int64("kafka.partition", int64(message.Partition)),
			attribute.Int64("kafka.offset", message.Offset),
		),
	)
	defer span.End()

	var orderResult oteldemo.OrderResult
	if err := proto.Unmarshal(message.Value, &orderResult); err != nil {
		log.Printf("Failed to unmarshal order: %v", err)
		span.RecordError(err)
		return
	}

	logOrderReceived(&orderResult)

	if c.db == nil {
		log.Println("Database connection not available, skipping persistence")
		return
	}

	if err := c.saveOrderToDB(&orderResult); err != nil {
		log.Printf("Failed to save order to database: %v", err)
		span.RecordError(err)
		return
	}

	span.SetAttributes(attribute.String("order.id", orderResult.OrderId))
	log.Printf("Successfully processed order: %s", orderResult.OrderId)
}

// saveOrderToDB saves the order to the database
func (c *Consumer) saveOrderToDB(orderResult *oteldemo.OrderResult) error {
	return c.db.Transaction(func(tx *gorm.DB) error {
		// Save order
		orderEntity := OrderEntity{
			ID: orderResult.OrderId,
		}
		if err := tx.Create(&orderEntity).Error; err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		// Save order items
		for _, item := range orderResult.Items {
			orderItem := OrderItemEntity{
				OrderID:              orderResult.OrderId,
				ProductID:            item.Item.ProductId,
				ItemCostCurrencyCode: item.Cost.CurrencyCode,
				ItemCostUnits:        item.Cost.Units,
				ItemCostNanos:        item.Cost.Nanos,
				Quantity:             item.Item.Quantity,
			}
			if err := tx.Create(&orderItem).Error; err != nil {
				return fmt.Errorf("failed to create order item: %w", err)
			}
		}

		// Save shipping info
		shipping := ShippingEntity{
			ShippingTrackingID:       orderResult.ShippingTrackingId,
			OrderID:                  orderResult.OrderId,
			ShippingCostCurrencyCode: orderResult.ShippingCost.CurrencyCode,
			ShippingCostUnits:        orderResult.ShippingCost.Units,
			ShippingCostNanos:        orderResult.ShippingCost.Nanos,
			StreetAddress:            orderResult.ShippingAddress.StreetAddress,
			City:                     orderResult.ShippingAddress.City,
			State:                    orderResult.ShippingAddress.State,
			Country:                  orderResult.ShippingAddress.Country,
			ZipCode:                  orderResult.ShippingAddress.ZipCode,
		}
		if err := tx.Create(&shipping).Error; err != nil {
			return fmt.Errorf("failed to create shipping: %w", err)
		}

		return nil
	})
}

// Close closes the consumer
func (c *Consumer) Close() error {
	return c.consumer.Close()
}

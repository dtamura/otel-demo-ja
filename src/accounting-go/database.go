// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupDatabase initializes the database connection
func setupDatabase() (*gorm.DB, error) {
	connectionString := os.Getenv("DB_CONNECTION_STRING")
	if connectionString == "" {
		log.Println("DB_CONNECTION_STRING not set, database operations will be skipped")
		return nil, nil
	}

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&OrderEntity{}, &OrderItemEntity{}, &ShippingEntity{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database schema: %w", err)
	}

	log.Println("Database connection established")
	return db, nil
}

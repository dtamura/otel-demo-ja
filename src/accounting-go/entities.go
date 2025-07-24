// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import "time"

// OrderEntity represents the order table
type OrderEntity struct {
	ID        string    `gorm:"column:order_id;primaryKey" json:"order_id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

// TableName sets the table name for OrderEntity
func (OrderEntity) TableName() string {
	return "order"
}

// OrderItemEntity represents the orderitem table
type OrderItemEntity struct {
	OrderID              string `gorm:"column:order_id;primaryKey" json:"order_id"`
	ProductID            string `gorm:"column:product_id;primaryKey" json:"product_id"`
	ItemCostCurrencyCode string `gorm:"column:item_cost_currency_code" json:"item_cost_currency_code"`
	ItemCostUnits        int64  `gorm:"column:item_cost_units" json:"item_cost_units"`
	ItemCostNanos        int32  `gorm:"column:item_cost_nanos" json:"item_cost_nanos"`
	Quantity             int32  `gorm:"column:quantity" json:"quantity"`
}

// TableName sets the table name for OrderItemEntity
func (OrderItemEntity) TableName() string {
	return "orderitem"
}

// ShippingEntity represents the shipping table
type ShippingEntity struct {
	ShippingTrackingID       string `gorm:"column:shipping_tracking_id;primaryKey" json:"shipping_tracking_id"`
	OrderID                  string `gorm:"column:order_id" json:"order_id"`
	ShippingCostCurrencyCode string `gorm:"column:shipping_cost_currency_code" json:"shipping_cost_currency_code"`
	ShippingCostUnits        int64  `gorm:"column:shipping_cost_units" json:"shipping_cost_units"`
	ShippingCostNanos        int32  `gorm:"column:shipping_cost_nanos" json:"shipping_cost_nanos"`
	StreetAddress            string `gorm:"column:street_address" json:"street_address"`
	City                     string `gorm:"column:city" json:"city"`
	State                    string `gorm:"column:state" json:"state"`
	Country                  string `gorm:"column:country" json:"country"`
	ZipCode                  string `gorm:"column:zip_code" json:"zip_code"`
}

// TableName sets the table name for ShippingEntity
func (ShippingEntity) TableName() string {
	return "shipping"
}

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"log"

	oteldemo "accounting-go/genproto"
)

// logOrderReceived logs the received order details
func logOrderReceived(orderResult *oteldemo.OrderResult) {
	orderJSON, err := json.Marshal(orderResult)
	if err != nil {
		log.Printf("Failed to marshal order result: %v", err)
		return
	}

	log.Printf("Order details: %s", string(orderJSON))
}

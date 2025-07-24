// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

var relevantPrefixes = []string{"GO_", "OTEL_", "KAFKA_", "DB_"}

// filterRelevantEnvVars filters environment variables by relevant prefixes
func filterRelevantEnvVars() map[string]string {
	envs := make(map[string]string)

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		for _, prefix := range relevantPrefixes {
			if strings.HasPrefix(strings.ToUpper(key), prefix) {
				envs[key] = value
				break
			}
		}
	}

	return envs
}

// outputEnvVarsInOrder outputs environment variables in alphabetical order
func outputEnvVarsInOrder(envs map[string]string) {
	var keys []string
	for key := range envs {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		fmt.Printf("%s=%s\n", key, envs[key])
	}
}

# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

FROM --platform=${BUILDPLATFORM} golang:1.22-alpine AS builder
ARG TARGETARCH
ARG TARGETOS
WORKDIR /app

# Install protobuf compiler
RUN apk add --no-cache protobuf-dev

# Install Go protobuf plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.0

# Copy proto file and generate Go code
COPY pb/demo.proto ./
RUN mkdir -p genproto
RUN protoc --go_out=./genproto --go_opt=paths=source_relative demo.proto

# Copy go mod files and download dependencies
COPY src/accounting-go/go.mod src/accounting-go/go.sum ./
RUN go mod download

# Copy source code
COPY src/accounting-go/*.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o accounting .

# -----------------------------------------------------------------------------

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/accounting .

# Create app user
RUN addgroup -g 1001 app && adduser -u 1001 -G app -s /bin/sh -D app
RUN chown app:app /app/accounting
USER app

ENTRYPOINT ["./accounting"]

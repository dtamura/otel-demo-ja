# Accounting Service (Go)

This service consumes new orders from a Kafka topic and is implemented in Go.

## Features

- Consumes orders from Kafka topic "orders"
- Stores order information in PostgreSQL database
- OpenTelemetry tracing support
- Structured logging
- Protobuf message deserialization

## Local Build

To build the service binary, run:

```sh
# From root context
mkdir -p src/accounting-go/genproto/
cp pb/demo.proto src/accounting-go/demo.proto

# In accounting-go service context
cd src/accounting-go
protoc --go_out=./genproto --go_opt=paths=source_relative demo.proto
go mod tidy
go build -o accounting .
```

## Docker Build

From the root directory, run:

```sh
docker compose build accounting-go
```

## Environment Variables

- `KAFKA_ADDR`: Kafka broker address (required)
- `DB_CONNECTION_STRING`: PostgreSQL connection string (optional, if not set, database operations are skipped)
- `OTEL_EXPORTER_OTLP_ENDPOINT`: OpenTelemetry collector endpoint (optional)

## Database Schema

The service creates and uses the following tables:

- `order`: Stores order information
- `orderitem`: Stores order line items
- `shipping`: Stores shipping information

## Dependencies

- Kafka client: IBM/sarama
- Database ORM: GORM with PostgreSQL driver
- Tracing: OpenTelemetry Go SDK
- Protobuf: google.golang.org/protobuf

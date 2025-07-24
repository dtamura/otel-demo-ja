# Accounting Service

This service consumes new orders from a Kafka topic.

**Note**: This service has been reimplemented in Go. See `../accounting-go/` for the Go version.

## Local Build

To build the C# service binary, run:

```sh
mkdir -p src/accounting/proto/ # root context
cp pb/demo.proto src/accounting/proto/demo.proto # root context
dotnet build # accounting service context
```

## Docker Build

From the root directory, run:

```sh
docker compose build accounting
```

## Bump dependencies

To bump all dependencies run in Package manager:

```sh
Update-Package -ProjectName Accounting
```

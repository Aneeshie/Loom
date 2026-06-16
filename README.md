# Loom

Loom is a distributed log aggregation system built with Go, gRPC, and PostgreSQL.

It provides:

* Log ingestion through lightweight agents
* Centralized log storage in PostgreSQL
* Historical log querying with filtering and pagination
* Realtime log streaming
* Realtime filtered subscriptions

---

## Architecture

```text
+--------+
| Agent  |
+--------+
     |
     | gRPC
     v
+----------------+
|   Loom Server  |
+----------------+
     |
     v
+------------+
| PostgreSQL |
+------------+

     ^
     |
     | gRPC
+--------+
| Query  |
+--------+
```

### Components

#### Agent

Reads logs from a source file, parses them, and sends them to the server using gRPC.

#### Server

Receives logs from agents, stores them in PostgreSQL, and streams logs to subscribed clients.

#### Query

CLI tool used to:

* Query historical logs
* Filter logs
* Follow logs in realtime

---

## Features

### Log Ingestion

Agents continuously read log files and send entries to the server.

### Historical Queries

Retrieve logs using filters such as:

* Log level
* Service name
* Host

Example:

```bash
go run ./cmd/query -level INFO
```

### Pagination

Limit the number of returned logs:

```bash
go run ./cmd/query -limit 20
```

### Realtime Streaming

Follow logs as they arrive:

```bash
go run ./cmd/query --follow
```

### Realtime Filtering

Subscribe only to logs matching specific criteria:

```bash
go run ./cmd/query --follow -level ERROR
```

---

## Technology Stack

* Go
* gRPC
* PostgreSQL
* pgx
* Protocol Buffers

---

## Running Loom

### Start PostgreSQL

Ensure PostgreSQL is running and accessible.

### Run Migrations

```bash
make migrate-up
```

### Start the Server

```bash
go run ./cmd/server
```

### Start an Agent

```bash
go run ./cmd/agent
```

### Query Historical Logs

```bash
go run ./cmd/query
```

Filter by level:

```bash
go run ./cmd/query -level INFO
```

Filter by service:

```bash
go run ./cmd/query -service auth
```

Limit results:

```bash
go run ./cmd/query -limit 10
```

### Follow Logs in Realtime

```bash
go run ./cmd/query --follow
```

Realtime filtered stream:

```bash
go run ./cmd/query --follow -level ERROR
```

---

## Testing

Run all tests:

```bash
go test ./...
```

Current test coverage includes:

* Log parsing
* Stream filtering
* Broadcast behavior

---

## Future Improvements

* Historical tail + follow mode
* Integration tests
* Structured log support
* Authentication and authorization
* Web dashboard
* Multiple storage backends

---

## Version

Current Release: **v1.0.0**

Loom v1.0.0 provides a complete end-to-end log aggregation pipeline with ingestion, storage, querying, filtering, and realtime streaming.


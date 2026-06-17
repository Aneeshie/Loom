# Loom

Loom is a distributed log aggregation system built with Go, gRPC, PostgreSQL, and pgvector.

It provides:

* Log ingestion through lightweight agents
* Centralized log storage in PostgreSQL
* Historical log querying with filtering
* Realtime log streaming
* Semantic log search using vector embeddings
* Hybrid search combining semantic retrieval and structured filters

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
     | Embeddings
     v
+------------------------+
| PostgreSQL + pgvector |
+------------------------+

     ^
     |
     | gRPC
+--------+
| Query  |
+--------+
```

---

## Components

### Agent

Reads logs from a source file, parses them, and sends them to the server using gRPC.

### Server

Receives logs from agents, generates vector embeddings, stores logs in PostgreSQL, and streams logs to subscribed clients.

### Query

CLI tool used to:

* Query historical logs
* Filter logs
* Follow logs in realtime
* Perform semantic search
* Perform hybrid search

---

## Features

### Log Ingestion

Agents continuously read log files and send entries to the server.

### Historical Queries

Retrieve logs using filters such as:

* Log level
* Service name
* Host
* Time range

```bash
go run ./cmd/query --level INFO
```

---

### Semantic Search

Search logs by meaning instead of exact keywords.

```bash
go run ./cmd/query --similar "database timeout"
```

Example matches:

```text
Database connection timeout
Query execution exceeded timeout
PostgreSQL unavailable
Connection pool exhausted
```

---

### Hybrid Search

Combine semantic search with structured filters.

```bash
go run ./cmd/query \
  --similar "database timeout" \
  --level ERROR
```

Example:

```text
Database connection timeout
Query execution exceeded timeout
PostgreSQL unavailable
```

---

### Realtime Streaming

Follow logs as they arrive.

```bash
go run ./cmd/query --follow
```

---

### Realtime Filtered Streaming

Subscribe only to logs matching specific filters.

```bash
go run ./cmd/query \
  --follow \
  --level ERROR
```

---

## Semantic Search Pipeline

```text
Query
 ↓
Embedding Generation
 ↓
Vector Similarity Search
 ↓
Structured Filtering
 ↓
Ranked Results
```

Logs are embedded using Ollama and stored as 768-dimensional vectors using pgvector.

---

## Technology Stack

* Go
* gRPC
* PostgreSQL
* pgvector
* Ollama
* Protocol Buffers
* pgx

---

## Running Loom

### Start PostgreSQL

Ensure PostgreSQL with pgvector support is running.

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

### Query Logs

All logs:

```bash
go run ./cmd/query
```

Filter by level:

```bash
go run ./cmd/query --level INFO
```

Semantic search:

```bash
go run ./cmd/query --similar "login success"
```

Hybrid search:

```bash
go run ./cmd/query \
  --similar "database timeout" \
  --level ERROR
```

Realtime stream:

```bash
go run ./cmd/query --follow
```

---

## Testing

Run all tests:

```bash
go test ./...
```

Current coverage includes:

* Log parsing
* Stream filtering
* Broadcast behavior
* Filter matching logic

---

## Roadmap

* Metadata extraction (JSONB enrichment)
* Similarity score exposure
* Natural language querying
* Historical tail + follow mode
* Integration tests
* Authentication and authorization
* Web dashboard

---

## Version

Current Release: **v1.2.0**

Loom v1.2.0 introduces semantic and hybrid log search powered by vector embeddings and pgvector.


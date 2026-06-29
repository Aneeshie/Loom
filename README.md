# Loom

Loom is a distributed log aggregation system built with Go, gRPC, PostgreSQL, and pgvector. It features a modern, interactive Next.js web dashboard and an AI-powered natural language query engine.

It provides:

* **Log Ingestion**: Lightweight agents tailing log sources and sending events via gRPC.
* **Centralized Storage**: Log records persisted securely in PostgreSQL.
* **Semantic Log Search**: Context-aware log searching powered by vector embeddings and pgvector.
* **Hybrid Search**: Seamlessly combines vector similarity search with structured metadata filters.
* **Intelligent Query Routing**: Bypasses the LLM parser when structured filters are provided, routing queries dynamically to optimize latency.
* **Web Dashboard**: Sleek, glassmorphic UI featuring log inspectors, interactive stats, and live telemetry filtering.

---

## System Architecture

```text
               +-----------------------+
               |   Web Dashboard UI    |
               |  (Next.js / Tailwind) |
               +-----------------------+
                           |
                           | HTTP / JSON API
                           v
+--------+     +-----------------------+
| Agent  |     |      Web Server       |
+--------+     |   (Go API Gateway)    |
     |         +-----------------------+
     |                     |
     | gRPC                | gRPC (LogService)
     v                     v
+--------------------------------------+
|             Loom Server              |
|        (Core Engine / Embeddings)    |
+--------------------------------------+
     |                     |
     | pgx                 | Ollama API
     v                     v
+------------------------+ +-----------+
| PostgreSQL + pgvector  | | Ollama LLM |
+------------------------+ +-----------+
```

---

## Components

### 1. Web Dashboard (Next.js)
A premium dashboard built with Next.js, Tailwind CSS, and Shadcn UI. It includes:
* **Interactive telemetry metrics** (filtered log count, real-time error rates, service/host statistics).
* **Advanced search bar** popover menu to filter logs by Severity, Service, Host, and Time window.
* **Logs viewer grid** featuring a sliding **Log Inspector Drawer** to view dynamic JSON context metadata and copy payloads.

### 2. Web Server (Go API Gateway)
Acts as the REST API gateway (`cmd/web`) serving `/api/v1/query`. It implements **Hybrid Query Routing**:
* **Natural Language Queries**: If a user submits a query string without any filters, it calls the LLM (Ollama) to extract fields (level, service, host, since).
* **Direct Semantic Search**: If query text AND filters are specified, it skips the LLM and calls `SimilarLogs` (vector search).
* **Direct Fetch**: If only filters are specified, it skips semantic search and LLM, performing a direct query (`GetLogs`).

### 3. Agent (Go)
Reads logs from configured file sources, parses severity and timing, and streams entries to the core server via gRPC.

### 4. Server (Go)
Accepts gRPC streams from agents, generates vector embeddings, stores logs, and manages subscriptions for real-time query clients.

---

## Running Loom

### 1. Start Database & Dependencies
Ensure PostgreSQL with `pgvector` and Ollama are running.
```bash
docker-compose up -d
```

### 2. Run Database Migrations
```bash
make migrate-up
```

### 3. Start Core Loom Server
```bash
go run ./cmd/server
```

### 4. Start Ingestion Agent
```bash
go run ./cmd/agent
```

### 5. Run the Go API Gateway (Web Server)
```bash
go run ./cmd/web
```
The Go server will start serving endpoints at `http://localhost:3000`.

### 6. Start the Web Dashboard
```bash
cd web
pnpm install
pnpm run dev
```
Open `http://localhost:3001` (or your terminal's allocated port) to access the interactive dashboard.

---

## API Query Endpoints

### Query Endpoint
* **Path**: `/api/v1/query`
* **Method**: `POST`
* **Payload**:
```json
{
  "query": "connection timeout",
  "level": "ERROR",
  "service": "payment-gateway",
  "host": "prod-ap-south-1a",
  "since": "1h"
}
```

---

## Roadmap

* [x] Web dashboard UI
* [x] Natural language querying integration (Ollama LLM)
* [x] Hybrid routing optimization
* [ ] Metadata extraction (JSONB enrichment)
* [ ] Similarity score exposure in UI
* [ ] Authentication and authorization

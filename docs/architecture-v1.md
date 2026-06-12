# Loom Architecture V1

## Overview

Loom is a self-hosted log aggregation platform.

Its purpose is to collect logs from multiple applications, centralize them, and provide a searchable interface for engineers.

Version 1 focuses only on ingestion, storage, and querying.

---

## Problem Statement

Modern applications generate logs across multiple services and machines.

When investigating failures, engineers often need to manually inspect log files from different systems.

Loom centralizes logs into a single system to simplify debugging and monitoring.

---

## Goals

Version 1 should support:

* Collecting logs from files.
* Sending logs to a central server.
* Persisting logs in PostgreSQL.
* Querying stored logs.

---

## Non-Goals

Version 1 does not include:

* Web UI
* WebSockets
* Semantic search
* pgvector
* Alerting
* Authentication
* Multi-tenancy
* LLM integration

---

## System Components

### Agent

Responsibilities:

* Watch configured log files.
* Read newly appended log lines.
* Send logs to the Loom Server.

Not responsible for:

* Storage
* Searching
* Alerting

---

### Server

Responsibilities:

* Receive logs from agents.
* Validate incoming data.
* Store logs in PostgreSQL.
* Provide APIs for querying logs.

Not responsible for:

* Reading files directly.

---

### PostgreSQL

Responsibilities:

* Persist log data.
* Serve as the source of truth.

---

## Data Flow

### Ingestion

Application
→ Log File
→ Loom Agent
→ gRPC
→ Loom Server
→ PostgreSQL

### Search

User
→ Query API
→ PostgreSQL
→ Results

---

## Repository Structure

```text
loom/

├── cmd/
│   ├── agent/
│   └── server/
│
├── internal/
│   ├── agent/
│   ├── ingestion/
│   ├── storage/
│   └── api/
│
├── proto/
├── migrations/
├── docs/
├── .github/
│   └── workflows/
│
├── docker-compose.yml
├── README.md
└── go.mod
```

---

## Future Versions

### V2

* Web UI

### V3

* Real-time log tailing (WebSockets)

### V4

* Semantic search using pgvector

### V5

* Alerting and anomaly detection


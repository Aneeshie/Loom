# Loom Architecture V2

## Overview
Loom V2 extends the distributed log aggregation architecture of V1 by adding a Go web server gateway, a modern Next.js dashboard client, and an AI-powered hybrid semantic query router.

---

## Extended System Components

### 1. Web Dashboard (Next.js App)
- **Role**: Serves as the developer portal for visualizing logs, viewing system metrics, and searching telemetry data.
- **Key Features**:
  - Stateful metrics calculations (Error rates, active servers, filtered query volume).
  - Floating search input with keybindings (`⌘K` / `/`).
  - Advanced popover menu supporting level, service, host, and duration parameters.
  - Interactive sliding drawer for inspecting full JSON log schemas.

### 2. Web API Gateway (`cmd/web` Server)
- **Role**: Serves as the HTTP boundary for frontend clients, hosting `/api/v1/query`.
- **Hybrid Routing Pipeline**:
  - To optimize latency and compute costs, the gateway inspects the incoming JSON payload and handles routing dynamically:
  
  ```text
  Is there any structured filter? (level, service, host, since)
         |
         +---> NO: Invoke LLM parser (Ollama) to extract fields from natural query.
         |
         +---> YES: Bypass LLM parser. Directly map payload parameters.
  
  After mapping parameters to Intent:
  Is there search query text?
         |
         +---> YES: Invoke similarLogs vector search (SimilarLogs gRPC client).
         |
         +---> NO: Invoke standard database query (GetLogs gRPC client).
  ```

### 3. Database (PostgreSQL + pgvector)
- **Role**: Stores raw logs and their corresponding 768-dimensional embeddings generated from Ollama.
- **Tables**:
  - `logs`: ID, service name, level, message, host, timestamp, and vector embeddings (`embedding` type).

---

## Ingestion Flow (Unchanged)
```text
Application Logs -> Agent -> gRPC (SendLog) -> Loom Server -> pgxpool -> PostgreSQL DB
```

## Search Query Flow (New V2)
```text
Next.js UI -> REST POST -> Go Web Gateway -> gRPC Client -> Loom Query Server -> similarity search (pgvector) / raw query -> JSON response
```

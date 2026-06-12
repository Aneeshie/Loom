# Architectural Decisions

This document records important architectural decisions made during development.

---

# ADR-001

## Decision

Use an installable agent instead of an application SDK.

## Status

Accepted

## Reasoning

An agent allows Loom to be language agnostic.

Applications do not need to import Loom-specific libraries or modify their code.

Any service capable of writing logs to a file can integrate with Loom.

---

# ADR-002

## Decision

Use PostgreSQL as the primary datastore.

## Status

Accepted

## Reasoning

PostgreSQL is operationally simple, reliable, and sufficient for Version 1.

It also provides a future upgrade path through pgvector for semantic search.

---

# ADR-003

## Decision

Use gRPC between Agent and Server.

## Status

Accepted

## Reasoning

gRPC provides:

* Strong contracts through Protocol Buffers.
* Efficient serialization.
* Streaming support for future versions.
* Cross-language compatibility.

---

# ADR-004

## Decision

Focus on ingestion and storage before advanced features.

## Status

Accepted

## Reasoning

A reliable ingestion pipeline is the foundation of the entire platform.

Features such as WebSockets, semantic search, and alerting will be built on top of a stable ingestion system.


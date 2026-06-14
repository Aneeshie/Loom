# Understanding Ownership, Dependencies, and Runtime Flow

This section exists because one of the most common backend confusions is:

> If Agent calls SendLog(), why does main create everything?

The answer is that there are **three different graphs** in every application.

Most confusion comes from accidentally mixing them together.

---

# Graph 1: Runtime Flow

This answers:

> Who talks to whom while the program is running?

For Loom:

```text
Agent
 ↓
gRPC
 ↓
LogService.SendLog()
 ↓
Store.InsertLog()
 ↓
Postgres
```

This is what happens when a log arrives.

Example:

```text
[INFO] User logged in
```

```text
Agent
 ↓
SendLog RPC
 ↓
LogService
 ↓
Store
 ↓
INSERT INTO logs
 ↓
Postgres
```

This graph is about execution.

It shows what happens after the application has already started.

---

# Graph 2: Dependency Graph

This answers:

> What does each component need in order to work?

For Loom:

```text
LogService
 ↓
needs
 ↓
Store

Store
 ↓
needs
 ↓
pgxpool.Pool

pgxpool.Pool
 ↓
needs
 ↓
Database URL
```

Notice:

```text
needs
```

does NOT mean:

```text
creates
```

This is where most beginners get confused.

---

# Graph 3: Ownership Graph

This answers:

> Who creates everything?

For Loom:

```text
main
 │
 ├── Config
 │
 ├── Store
 │
 ├── LogService
 │
 └── GRPCServer
```

Everything starts here.

Think of main as:

```text
The assembly room
```

Nothing important gets created anywhere else.

---

# The Restaurant Analogy

Imagine Loom is a restaurant.

Runtime flow:

```text
Customer
 ↓
Waiter
 ↓
Chef
 ↓
Kitchen
```

Ownership:

```text
Restaurant Owner
 ├── Waiter
 ├── Chef
 └── Kitchen
```

Question:

Does the waiter hire the chef?

No.

The waiter uses the chef.

The owner hires the chef.

---

Loom is the same.

```text
Agent
 ↓
LogService
 ↓
Store
```

LogService uses Store.

But LogService does not create Store.

main creates Store.

---

# Why main Creates Store

Suppose we write:

```go
func NewLogService() *LogService {
    store := storage.NewStore(...)
    return &LogService{
        store: store,
    }
}
```

Now LogService secretly owns Store.

This causes problems:

```text
Can't replace Store
Can't test easily
Hidden dependency
```

Instead:

```go
store := storage.NewStore(...)

logService := grpc.NewLogService(store)
```

Now:

```text
main owns Store

LogService uses Store
```

This is dependency injection.

---

# Startup Phase vs Runtime Phase

This distinction is the most important idea.

Startup:

```text
main
 ↓
Load Config
 ↓
Create Store
 ↓
Create LogService
 ↓
Create GRPCServer
 ↓
Start Server
```

Runtime:

```text
Agent
 ↓
SendLog()
 ↓
LogService
 ↓
Store
 ↓
Postgres
```

Startup happens once.

Runtime happens thousands of times.

---

# The Exact Loom Flow

Server startup:

```go
cfg := LoadServerConfig()

store := storage.NewStore(cfg.Database.URL)

logService := grpc.NewLogService(store)

server := grpc.NewServer(logService)

server.Run()
```

Hours later:

```text
Agent sends log
```

gRPC calls:

```go
logService.SendLog(...)
```

Inside SendLog:

```go
s.store.InsertLog(...)
```

Why does this work?

Because the store was attached to the LogService during startup.

The LogService remembers it.

---

# Mental Model

When designing a component, ask three questions:

1. Runtime

```text
Who calls me?
```

2. Dependency

```text
What do I need?
```

3. Ownership

```text
Who creates me?
```

For LogService:

```text
Who calls me?
→ gRPC

What do I need?
→ Store

Who creates me?
→ main
```

For Store:

```text
Who calls me?
→ LogService

What do I need?
→ pgxpool

Who creates me?
→ main
```

Answering those three questions is usually enough to design the entire architecture.

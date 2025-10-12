# 🚦 Thread-Safe HTTP Load Balancer

A custom-built HTTP load balancer designed with **backend agnosticism**, **session persistence**, and **thread-safe server rotation**.

## ✅ Feature Checklist

- ✅ **Systems-level concurrency control**
- ✅ **Custom data structures (doubly linked list round-robin for server scheduling)**
- ✅ **Cookie-based load balancing strategy**
- ✅ **Integration & unit tests to guarantee thread safety**
- ✅ **Redis-backed distributed session storage** (in progress)

---

## 🧠 Design Overview

### Load Balancing Strategy

| Component | Purpose |
|-----------|---------|
| **Doubly Linked List (Round Robin)** | Maintains the active server pool. Each request rotates through nodes in a lock-safe manner. |
| **Injected Session Cookie** | The first response assigns a backend ID to the client. Subsequent requests can be routed to any backend server as the redis allows for agnosticism. |
| **Redis Session Store** | Holds all user session data centrally so backends remain stateless and interchangeable. |
| **Thread-Safety Guarantees** | The server list is stress-tested under concurrency to ensure no race conditions  |

---

## Deployment

Coming Soon


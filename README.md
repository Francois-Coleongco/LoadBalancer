# 🚦 Thread-Safe HTTP Load Balancer

A custom-built HTTP load balancer designed with **backend agnosticism**, **session persistence**, and **thread-safe server rotation**.

Server Rotation Demo:

https://github.com/user-attachments/assets/96132742-b8ca-4fa6-8886-f9c66e4dd74c



## ✅ Feature Checklist

- ✅ **Systems-level concurrency control**
- ✅ **Custom data structures (doubly circular linked list round-robin for server scheduling)**
- ✅ **Cookie-based load balancing strategy**
- ✅ **Functional Tests to guarantee thread safety and validated server insertion, deletion, and traversal edge cases (single node, full rotation, and empty state)**
- ✅ **Redis-backed distributed session storage** (in progress)

---

## 🧠 Design Overview

### Load Balancing Strategy

| Component | Purpose |
|-----------|---------|
| **Doubly Circular Linked List (Round Robin)** | Maintains the active server pool. Each request rotates through nodes in a lock-safe manner. |
| **Injected Session Cookie** | The first response assigns a backend ID to the client. Subsequent requests can be routed to any backend server as the redis allows for agnosticism. |
| **Redis Session Store** | Holds all user session data centrally so backends remain stateless and interchangeable. |
| **Thread-Safety Guarantees** | The server list is stress-tested under concurrency to ensure no race conditions  |

---

## Deployment

Coming Soon


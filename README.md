# 🚦 Thread-Safe HTTP Load Balancer

<img width="2387" height="939" alt="system_design_diagram" src="https://github.com/user-attachments/assets/f8e3d989-c0ac-4eb9-85c7-5ad2778581a2" />

A custom-built HTTP load balancer designed with **backend agnosticism**, **session persistence**, and **thread-safe server rotation**.

Server Rotation + Redis HTTP Header Storage Demo:

https://github.com/user-attachments/assets/a90ebaea-c2db-4007-9751-1439ab60e0de


## ✅ Feature Checklist

- ✅ **Systems-level concurrency control**
- ✅ **Custom data structures (doubly circular linked list round-robin for server scheduling)**
- ✅ **Cookie-based load balancing strategy**
- ✅ **Functional Tests to guarantee thread safety and validated server insertion, deletion, and traversal edge cases (single node, full rotation, and empty state)**
- ✅ **Redis-backed distributed session storage**

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

## Environment Configuration & Building

### Server File:

setup your ports like so in a file:

🔧 `servers.txt`

```
http://localhost:8000
http://localhost:8001
http://localhost:8080
http://localhost:8081
http://localhost:5000
http://localhost:3000
http://localhost:8888
http://localhost:5173
```

## Usage


Start the Redis store:

```docker run -d -p 6379:6379 --name lb_redis_store redis```

The load balancer takes the following flags:

-f servers.txt
-p loadBalancerPort

The full command looks something like:

```./main -f servers.txt -p 7777```


### Notes:

If your backends use JWTs or another stateless session mechanism, you likely don’t need Redis for session routing. The session information is contained in the JWT and sent with every request, so any backend can handle it without sticky sessions.

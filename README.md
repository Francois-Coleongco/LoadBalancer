# Thread-Safe HTTP Load Balancer

<img width="2270" height="978" alt="system_design_diagram" src="https://github.com/user-attachments/assets/c98ac8f7-789d-42d2-a738-452c75e59b7a" />

## Server Rotation Demo:

https://github.com/user-attachments/assets/f8307562-1a2b-4d7f-ad39-388eeb176a89



A custom-built HTTP load balancer designed with **backend agnosticism**, **session persistence**, and **thread-safe server rotation** for distributed systems.


## Feature list

- **Systems-level concurrency control**
- **Custom data structures (doubly circular linked list round-robin for server scheduling)**
- **Cookie-based load balancing strategy**
- **Functional Tests to guarantee thread safety and validated server insertion, deletion, and traversal edge cases (single node, full rotation, and empty state)**
- **Redis-backed session storage**

---

## Design Overview

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

Setup your ports like so in a file with a name of your choice:

 `servers.txt`

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


Start the Redis store and Redis CLI (this is what your backend will talk to in order to retrieve session data):

```
docker run -d -p 6379:6379 --name lb_redis_store redis
docker exec -it lb_redis_store redis-cli
```

The load balancer takes the following flags:

```
-f your_servers_file.txt
-p your_load_balancer_port
```

The full command looks something like:

```
./main -f servers.txt -p 7777
```

## FUTURE:

I'm hoping to work on a better algorithm than round-robin, perhaps using the standard deviation of each server with the mean of connections to calculate the best server to route to next.

### Notes:

If your backends use JWTs or another stateless session mechanism, you likely don’t need and probably SHOULD NOT use Redis for session routing. The session information is contained in the JWT and sent with every request, so any backend can handle it without sticky sessions. Using the redis store in this case would be a lapse in security as the JWT would be stored unnecessarily for the duration of the redis (key, value) expiration.

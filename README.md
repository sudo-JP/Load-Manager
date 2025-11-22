# Load-Manager

This project is my little sandbox for playing with distributed systems.  
It’s built entirely in **Go**, partly to learn the language and partly to explore how load distribution, backend scaling, and system stress testing all fit together.

## Project Overview

The project is organized into three independent Go modules:

- `backend/`: Basic Go HTTP server (simulates a small service w/ Postgres)
- `load-manager/`: Load balancer / scheduler that routes traffic to backends
- `test/`: Stress tester that floods the system to measure performance

Also I tried some cryptography in `backend/internal/hash`

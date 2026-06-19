# ShardLab

ShardLab is a Go CLI and HTTP API for simulating distributed key sharding before changing cache or storage topology.

It compares modulo hashing, consistent hashing, and rendezvous hashing under node additions, node removals, virtual-node settings, weighted nodes, and hot-key workloads. The goal is to answer a practical distributed-systems question before production traffic moves:

> If I add or remove a node, how much traffic moves and how balanced will the cluster be afterward?

## Why This Project Exists

Backend teams often learn the hard way that a naive sharding function can move most keys during scaling events. ShardLab makes that failure mode visible with deterministic simulations and metrics.

It is intentionally not a CRUD app. It demonstrates:

- deterministic distributed-systems simulation
- modulo, consistent-hash-ring, and rendezvous assignment
- virtual nodes and node weighting
- hot-key workload modeling
- movement ratio, load skew, Gini coefficient, and overload flags
- CLI reports, JSON output, HTTP API, and SSE progress streaming
- tests, samples, Dockerfile, and CI

## Quick Start

```powershell
go test ./...
go run ./cmd/shardlab -input samples/add-node-consistent.json
go run ./cmd/shardlab -input samples/add-node-modulo.json
```

Write a JSON report:

```powershell
go run ./cmd/shardlab -input samples/hot-key-consistent.json -format json -out reports/hot-key-report.json
```

Run the API:

```powershell
go run ./cmd/api -addr :8080
```

Simulate with the API:

```powershell
Invoke-RestMethod -Method Post -Uri http://localhost:8080/api/simulate -ContentType application/json -InFile samples/add-node-consistent.json
```

Stream progress events:

```powershell
curl.exe -N -X POST http://localhost:8080/api/simulate/stream -H "Content-Type: application/json" --data-binary "@samples/add-node-consistent.json"
```

## API

```text
GET  /health
POST /api/simulate
POST /api/simulate/stream
```

## Metrics

| Metric | Meaning |
| --- | --- |
| `movementRatio` | Share of requests assigned to a different node after topology change |
| `loadSkew` | Max node load divided by mean node load |
| `giniCoefficient` | Inequality score across node loads |
| `overloadedNodes` | Nodes above 125% of mean load |
| `underloadedNodes` | Nodes below 75% of mean load |

## Why It Is Useful In Interviews

ShardLab gives interviewers concrete backend/distributed-systems questions:

- Why does modulo hashing move many keys when a node is added?
- How do virtual nodes improve consistent hashing?
- When would rendezvous hashing be simpler than a hash ring?
- How do hot keys break otherwise balanced partitioning?
- What metrics should a backend team check before scaling a cache cluster?

## Tech

Go 1.26, stdlib HTTP, deterministic SHA-256-based hashing, CLI, SSE streaming, tests, Docker, GitHub Actions.

# Testing Guide

Run all tests:

```powershell
go test ./...
```

Run the CLI:

```powershell
go run ./cmd/shardlab -input samples/add-node-consistent.json
go run ./cmd/shardlab -input samples/add-node-modulo.json
go run ./cmd/shardlab -input samples/hot-key-consistent.json
```

Expected behavior:

- modulo hashing should report a high movement ratio after adding a node
- consistent hashing should move a smaller share of requests
- hot-key workloads should show higher load skew

Run the API:

```powershell
go run ./cmd/api -addr :8080
```

Smoke test:

```powershell
Invoke-RestMethod -Uri http://localhost:8080/health
Invoke-RestMethod -Method Post -Uri http://localhost:8080/api/simulate -ContentType application/json -InFile samples/add-node-consistent.json
curl.exe -N -X POST http://localhost:8080/api/simulate/stream -H "Content-Type: application/json" --data-binary "@samples/add-node-consistent.json"
```

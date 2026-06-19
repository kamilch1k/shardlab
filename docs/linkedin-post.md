# LinkedIn Post Draft

I built ShardLab, a Go backend project for simulating distributed key sharding before changing cache or storage topology.

It compares modulo hashing, consistent hashing, and rendezvous hashing under node additions and hot-key workloads. The API reports key/request movement, load skew, Gini coefficient, overloaded nodes, and a recommendation.

The goal was to make a backend portfolio project that is both practical and a little more scientific:

- deterministic distributed-systems simulation
- hash-ring and rendezvous assignment algorithms
- hot-key workload modeling
- CLI + HTTP API + SSE progress stream
- tests, samples, Docker, and CI

Repo: https://github.com/kamilch1k/shardlab

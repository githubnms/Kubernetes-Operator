# Kubernetes Custom Metrics Autoscaler

> Status: In Progress — just started. This README is updated as the
> project develops. See `PROGRESS.md` for the detailed build log and
> `PROBLEM_STATEMENT.md` for the full scope and motivation.

## What This Is

A Kubernetes Operator, written in Go, that scales application pods based on
a custom-defined metric — instead of only CPU/memory like the default
Kubernetes Horizontal Pod Autoscaler.

## Why This Project

*(To be filled in as we build — short version: to learn Go and Kubernetes
internals by building a real controller, not just using existing tools.
Full motivation in `PROBLEM_STATEMENT.md`.)*

## Architecture

*(Diagram and explanation to be added once Phase 3 is underway. Will show:
`ScalingPolicy` CRD → Controller → Metric Source → Target Deployment.)*

## Tech Stack

- Go
- Kubernetes (`client-go`, `controller-runtime`)
- Docker
- Prometheus (planned, for metrics exposure)
- Minikube / Kind (local test cluster)

## Project Structure

```
Kubernetes-Operator/
├── cmd/
│   └── manager/              # main.go — entrypoint for the operator binary
├── api/
│   └── v1alpha1/             # CRD type definitions (ScalingPolicy)
├── controllers/              # Reconciliation loop logic
├── config/
│   ├── crd/                  # Generated CRD YAML manifests
│   ├── rbac/                 # RBAC permissions for the operator
│   ├── manager/              # Deployment manifest for the operator itself
│   └── samples/              # Example ScalingPolicy YAML resources
├── docs/                     # Diagrams, design notes
├── test/                     # Test scenarios / load-spike simulation scripts
├── hack/                     # Dev scripts (local cluster setup, etc.)
├── Dockerfile
├── Makefile
├── go.mod
├── PROBLEM_STATEMENT.md      # What problem this solves and why
├── PROGRESS.md               # Phased build log, updated as we go
└── README.md                 # This file
```

## How to Run

*(To be filled in once there's actually something runnable — Phase 3+.)*

## How to Test

*(To be filled in once test scenarios exist — Phase 6.)*

## Results / Numbers

*(To be filled in with real, measured numbers once testing is done — see
`PROGRESS.md` log entries for raw data as it's collected.)*

## What I Learned

*(To be filled in at the end — design trade-offs, what was hard, what would
be done differently.)*

-----

See `PROBLEM_STATEMENT.md` for the full list.
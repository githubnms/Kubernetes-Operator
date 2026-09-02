# Problem Statement

## The Problem

> Kubernetes' built-in Horizontal Pod Autoscaler (HPA) scales workloads based on
> CPU and memory usage by default. In many real systems, CPU/memory is a poor
> proxy for actual load. Examples:

- A worker service might be CPU-idle but have a growing backlog in a queue.
- An API service might have low CPU usage but degrading latency under
  connection load.
- A batch service might need to scale based on "items waiting to be
  processed," not resource usage at all.

> Teams solve this today with the Kubernetes Custom Metrics API and tools like
> KEDA — but understanding *how* that mechanism works internally (reconciliation
> loops, custom resources, controller patterns) is a distinct skill from being
> able to use an existing tool.

## Goal of This Project

Build a Kubernetes Operator, from scratch, in Go, that:

1. Introduces a Custom Resource Definition (CRD) called `ScalingPolicy`,
   where a user declares: which Deployment to scale, which metric to watch,
   and the thresholds/bounds for scaling.
2. Runs a controller that continuously reconciles the actual state of the
   cluster against the desired state declared in `ScalingPolicy` objects.
3. Reads a metric value (starting with a simple/mock metric source, later a
   real one such as Prometheus or an HTTP endpoint).
4. Scales the target Deployment's replica count up or down based on that
   metric, safely (respecting min/max bounds, avoiding rapid flapping).
5. Is observable — exposes logs and metrics about its own scaling decisions.

## Why This Project (Personal Motivation)

To learn Go by building something with real engineering constraints
(concurrency, eventual consistency, failure handling) rather than isolated
tutorial exercises, and to end up with a portfolio piece that demonstrates
understanding of the same patterns used by Kubernetes' own ecosystem
(controller-runtime, client-go) — technology used internally at companies
like Google, Amazon, Microsoft, and in the CNCF ecosystem generally (Docker,
Kubernetes itself, Prometheus).

## Non-Goals (Explicitly Out of Scope, at Least for v1)

To keep this achievable and honest about scope, this project will **not**,
at least initially:

- Run in a real production cluster or handle real production traffic.
- Support every metric source (starting with one: a mock or Prometheus).
- Handle multi-cluster or multi-region scaling.
- Replace or compete with KEDA/HPA — it's a learning re-implementation of
  the underlying pattern, not a production alternative.

Non-goals may be revisited later in `PROGRESS.md` as the project matures.

## Definition of Done (v1)

- [ ] `ScalingPolicy` CRD is defined and installable on a cluster.
- [ ] Controller watches `ScalingPolicy` objects and reconciles them.
- [ ] Controller reads a metric value from a configurable source.
- [ ] Controller scales the target Deployment's `replicas` field accordingly.
- [ ] Scaling respects `minReplicas` / `maxReplicas` bounds.
- [ ] Basic cooldown logic prevents rapid scale up/down flapping.
- [ ] Unit tests cover the reconciliation logic.
- [ ] Project runs successfully on a local Minikube/Kind cluster.
- [ ] README documents architecture, setup, and how to test it.
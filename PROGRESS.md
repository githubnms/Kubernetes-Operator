# Progress Log

**Status: In Progress — just started.**

---

## Phase 1 — Bootstrap Project

- [x] Initialize Go module (`go.mod`)
- [x] Scaffold project structure (cmd/, api/, controllers/, config/)
- [x] Add Dockerfile and Makefile skeleton

## Phase 2 — Define the CRD

- [x] Define `ScalingPolicy` custom resource type (`api/v1alpha1`)
- [ ] Fields: target Deployment reference, metric source, threshold,
      minReplicas, maxReplicas
- [ ] Generate CRD YAML manifest (`config/crd`)

## Phase 3 — Build the Reconciliation Loop

- [ ] Scaffold controller (`controllers/scalingpolicy_controller.go`)
- [ ] Watch `ScalingPolicy` objects and trigger reconcile on change
- [ ] Fetch current state of target Deployment
- [ ] Read metric value from a mock/configurable source
- [ ] Compare metric against threshold, decide scale up/down/no-op
- [ ] Apply new replica count to target Deployment
- [ ] Enforce `minReplicas` / `maxReplicas` bounds
- [ ] Add cooldown logic to prevent rapid scale flapping

## Phase 4 — Observability & Failure Handling

- [ ] Add structured logging for every reconciliation decision
- [ ] Expose Prometheus metrics (scale events, current replica count,
      last observed metric value)
- [ ] Handle failure cases: metric source unreachable, target Deployment
      missing, conflicting updates
- [ ] Unit tests for reconciliation logic

## Phase 5 — Package & Deploy

- [ ] Build Docker image for the operator
- [ ] Deploy CRD + operator to a local cluster (Minikube/Kind)
- [ ] Verify end-to-end: apply a `ScalingPolicy`, simulate metric change,
      confirm pod count changes
- [ ] (Stretch) Helm chart for deployment

## Phase 6 — Test, Measure, Document

- [ ] Run simulated load-spike scenarios, record real numbers:
      - Number of test runs: ___
      - Scaling response time observed: ___
      - Failures encountered and how handled: ___
- [ ] Finalize README (architecture diagram, setup, results)
- [ ] Write LinkedIn article + short post
- [ ] Finalize resume bullets using real numbers from testing

----
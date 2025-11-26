# Sluice (TCO Configurator)

> A Kubernetes-native log budget management system for controlling cloud observability costs

Log volumes in cloud-native systems can grow uncontrollably, leading to unpredictable storage and processing costs. TCO Configurator provides automated budget enforcement at the infrastructure level, allowing teams to maintain observability while controlling costs.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Design Patterns](#design-patterns)
- [Features](#features)
- [Quick Start](#quick-start)
- [Deployment](#deployment)
- [Development](#development)
- [API Reference](#api-reference)
- [License](#license)

## Overview

TCO Configurator integrates with Promtail, Kubernetes, and Prometheus to:
- **Track** log volume per app/team/namespace in real-time
- **Enforce** budget limits via dynamic log filtering
- **Provide** real-time visibility into cost vs. usage
- **Help** teams control costs without sacrificing observability

### Problem Statement

In cloud-native environments, log volumes can spike unexpectedly due to:
- Application errors generating excessive logs
- Debug logging left enabled in production
- Chatty microservices
- Lack of visibility into per-team costs

This leads to:
- Unpredictable monthly bills
- Budget overruns
- Difficulty attributing costs to teams
- No automated enforcement mechanism

### Solution

TCO Configurator provides:
1. **Declarative Budget Management** - Define budgets as Kubernetes resources
2. **Automated Enforcement** - Real-time log filtering based on usage
3. **Team Isolation** - Per-namespace budget tracking
4. **Graceful Degradation** - Throttle before dropping logs
5. **Observability** - Prometheus metrics for cost tracking

## Architecture

### System Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Kubernetes Cluster                          │
│                                                                     │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │                    Control Plane                             │   │
│  │                                                              │   │
│  │  ┌────────────────┐         ┌──────────────────┐             │   │
│  │  │  Controller    │◄────────│  TeamBudget CRD  │             │   │
│  │  │  (Deployment)  │  Watch  │   (API Server)   │             │   │
│  │  └────────┬───────┘         └──────────────────┘             │   │
│  │           │                                                  │   │
│  │           │ Reconcile                                        │   │
│  │           ▼                                                  │   │
│  │  ┌─────────────────────────────────────────────────────┐     │   │
│  │  │           Policy Engine (Stateless)                 │     │   │
│  │  │  ┌──────────────────────────────────────────────┐   │     │   │
│  │  │  │  EvaluateUsageStateless()                    │   │     │   │
│  │  │  │  - Input: team, limit, current, new          │   │     │   │
│  │  │  │  - Output: Allow/Throttle/Drop               │   │     │   │
│  │  │  │  - Logic: Stateless, no side effects         │   │     │   │
│  │  │  └──────────────────────────────────────────────┘   │     │   │
│  │  └─────────────────────────────────────────────────────┘     │   │
│  └──────────────────────────────────────────────────────────────┘   │
│                                                                     │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │                      Data Plane                              │   │
│  │                                                              │   │
│  │  ┌────────────────┐         ┌──────────────────┐             │   │
│  │  │  Agent         │◄────────│  Promtail Plugin │             │   │
│  │  │  (DaemonSet)   │  Query  │  (Log Pipeline)  │             │   │
│  │  └────────┬───────┘         └──────────────────┘             │   │
│  │           │                                                  │   │
│  │           │ ProcessLogs()                                    │   │
│  │           ▼                                                  │   │
│  │  ┌─────────────────────────────────────────────────────┐     │   │
│  │  │  1. Get TeamBudget from K8s API                     │     │   │
│  │  │  2. Call Policy Engine                              │     │   │
│  │  │  3. Update TeamBudget Status                        │     │   │ 
│  │  │  4. Return Action to Promtail                       │     │   │ 
│  │  └─────────────────────────────────────────────────────┘     │   │
│  └──────────────────────────────────────────────────────────────┘   │
│                                                                     │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │                   Observability Layer                        │   │
│  │                                                              │   │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────────────┐      │   │
│  │  │ Dashboard  │  │ Prometheus │  │  Notifier          │      │   │
│  │  │ (Web UI)   │  │ (Metrics)  │  │  (Alerts)          │      │   │
│  │  └────────────┘  └────────────┘  └────────────────────┘      │   │
│  └──────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
```

### Component Interaction Flow

```
┌─────────┐                                                    ┌──────────┐
│  App    │                                                    │   Loki   │
│  Logs   │                                                    │ (Storage)│
└────┬────┘                                                    └─────▲────┘
     │                                                               │
     │ 1. Log Entry                                                  │
     ▼                                                               │
┌─────────────┐                                                      │
│  Promtail   │                                                      │
│   Plugin    │                                                      │
└──────┬──────┘                                                      │
       │                                                             │
       │ 2. Handle(namespace, logSize)                               │
       ▼                                                             │
┌─────────────┐                                                      │
│    Agent    │                                                      │
│ ProcessLogs │                                                      │
└──────┬──────┘                                                      │
       │                                                             │
       │ 3. GetTeamBudget(namespace)                                 │
       ▼                                                             │
┌─────────────┐                                                      │
│ K8s API     │                                                      │
│ TeamBudget  │                                                      │
└──────┬──────┘                                                      │
       │                                                             │
       │ 4. Return: dailyLimit, currentUsage                         │
       ▼                                                             │
┌─────────────────────────────────┐                                  │
│  Policy Engine                  │                                  │
│  EvaluateUsageStateless()       │                                  │
│  - Calculate: total = current + new                                │
│  - Decide: Allow/Throttle/Drop  │                                  │
└──────┬──────────────────────────┘                                  │
       │                                                             │
       │ 5. Return: Action                                           │
       ▼                                                             │
┌─────────────┐                                                      │
│    Agent    │                                                      │
│ UpdateStatus│                                                      │
└──────┬──────┘                                                      │
       │                                                             │
       │ 6. UpdateTeamBudgetStatus(newUsage)                         │
       ▼                                                             │
┌─────────────┐                                                      │
│ K8s API     │                                                      │
│ TeamBudget  │                                                      │
└──────┬──────┘                                                      │
       │                                                             │
       │ 7. Status Updated                                           │
       ▼                                                             │
┌─────────────┐                                                      │
│  Promtail   │                                                      │
│   Plugin    │                                                      │
└──────┬──────┘                                                      │
       │                                                             │
       │ 8. If Allow: Forward Log ───────────────────────────────────┘
       │    If Throttle: Sample Log
       │    If Drop: Discard Log
       └─────────────────────────────────────────────────────────────
```

### Data Flow

```
TeamBudget CRD (Desired State)
        │
        │ Watch
        ▼
   Controller ──────► Policy Engine (Add to memory)
        │
        │ Reconcile Loop
        └──────────────────────────────────────┐
                                               │
Log Entry → Promtail → Agent → Policy Engine   │
                         │           │         │
                         │           ▼         │
                         │      Evaluate       │
                         │           │         │
                         │           ▼         │
                         │      Action         │
                         │           │         │
                         ▼           ▼         ▼
                    Update Status → K8s API ← Controller
                         │
                         ▼
                    Prometheus Metrics
```

## Design Patterns

### 1. Operator Pattern

**Why:** Kubernetes-native way to manage custom resources

**Implementation:**
- Custom Resource Definition (CRD) for `TeamBudget`
- Controller watches for CRD changes
- Reconciliation loop ensures desired state

**Benefits:**
- Declarative configuration
- GitOps compatible
- Native K8s integration
- Automatic state management

```go
// Controller reconciles TeamBudget resources
type Controller struct {
    KubeClient   *Client
    PolicyEngine *policy.PolicyEngine
}

func (c *Controller) ReconcileTeambudget(ctx context.Context, name, namespace string) error {
    // 1. Get TeamBudget from K8s
    // 2. Convert to internal Budget type
    // 3. Add to PolicyEngine
    // 4. Evaluate current state
}
```

### 2. Stateless Policy Engine

**Why:** Enables horizontal scaling and simplifies testing

**Implementation:**
- No internal state storage
- All data passed as parameters
- Pure function evaluation
- State stored in K8s CRD status

**Benefits:**
- Easy to test (no mocks needed)
- Horizontally scalable
- No state synchronization issues
- Crash-safe (state in K8s)

```go
// Stateless evaluation - no side effects
func EvaluateUsageStateless(team string, dailyLimit, currentUsage, newBytes int64) Action {
    totalBytes := currentUsage + newBytes
    
    switch {
    case totalBytes > dailyLimit:
        return ActionDrop
    case totalBytes > dailyLimit*80/100:
        return ActionThrottle
    default:
        return ActionAllow
    }
}
```

**Why Stateless?**
1. **Reliability:** If pod crashes, no state is lost
2. **Scalability:** Can run multiple instances without coordination
3. **Testability:** Pure functions are easy to test
4. **Simplicity:** No need for state synchronization

### 3. Interface-Based Design

**Why:** Enables testing without real dependencies

**Implementation:**
- Define interfaces for external dependencies
- Use dependency injection
- Mock implementations for tests

**Benefits:**
- Unit tests don't need K8s cluster
- Easy to swap implementations
- Loose coupling
- Better testability

```go
// Interface allows mocking in tests
type KubeClientInterface interface {
    GetTeamBudget(name, namespace string) (*v1.TeamBudget, error)
    UpdateTeamBudgetStatus(name, namespace string, newUsage int64) error
}

// Agent depends on interface, not concrete type
type Agent struct {
    KubeClient KubeClientInterface  // Can be real or mock
}

// Mock for testing
type MockKubeClient struct {
    dailyLimit   int64
    currentUsage int64
}
```

### 4. DaemonSet Pattern

**Why:** Agent needs to run on every node

**Implementation:**
- Agent deployed as DaemonSet
- One pod per node
- Intercepts logs at source

**Benefits:**
- Low latency (local processing)
- No single point of failure
- Scales with cluster
- Minimal network overhead

### 5. Watch Pattern

**Why:** React to changes in real-time

**Implementation:**
- Controller watches TeamBudget resources
- Event-driven reconciliation
- Automatic updates

**Benefits:**
- Real-time updates
- No polling overhead
- Efficient resource usage
- Immediate response to changes

```go
func (c *Controller) Start(ctx context.Context, namespace string) error {
    for {
        watcher, err := watchRes.Watch(ctx, metav1.ListOptions{})
        
        for event := range watcher.ResultChan() {
            if event.Type == watch.Added || event.Type == watch.Modified {
                // Reconcile the resource
                c.ReconcileTeambudget(ctx, name, ns)
            }
        }
    }
}
```

### 6. Graceful Degradation

**Why:** Don't drop logs immediately when over budget

**Implementation:**
- Three-tier action system
- Progressive enforcement
- Throttle before drop

**Benefits:**
- Better user experience
- Time to react
- Maintains some observability
- Reduces surprise

```
Usage Level          Action        Behavior
─────────────────────────────────────────────
0-80% of limit      Allow         Pass all logs
80-100% of limit    Throttle      Sample logs (e.g., 50%)
>100% of limit      Drop          Discard logs
```

### 7. Metrics-Driven

**Why:** Observability into the system itself

**Implementation:**
- Prometheus metrics for all actions
- Counter per team per action
- Exportable via /metrics endpoint

**Benefits:**
- Visibility into enforcement
- Alerting on budget violations
- Historical analysis
- Cost attribution

```go
// Metrics tracked per team and action
var TeamBudgetActionsCounter = promauto.NewCounterVec(
    prometheus.CounterOpts{
        Name: "tco_policy_engine_actions_total",
        Help: "Actions taken by policy engine",
    },
    []string{"team", "action"},
)
```

## Features

### Core Features

- **Log Volume Tracking** - Monitor usage per app/team/namespace in real-time
- **Budget Enforcement** - Automatically drop/throttle logs when over budget
- **Prometheus Metrics** - Real-time cost dashboards and alerting
- **K8s Native** - CRDs for team-level budget policies
- **Graceful Degradation** - Throttle before dropping logs
- **Multi-Tenant** - Per-namespace isolation and tracking

### Technical Features

- **Stateless Architecture** - Horizontally scalable, crash-safe
- **Event-Driven** - Real-time response to budget changes
- **Interface-Based** - Fully testable without K8s cluster
- **Prometheus Integration** - Built-in metrics export
- **RBAC Support** - Proper Kubernetes security
- **In-Cluster Auth** - Works with K8s service accounts

## Quick Start

### Prerequisites

- Kubernetes cluster (1.20+)
- kubectl configured
- Docker (for building images)

### 1. Install the CRD

```bash
kubectl apply -f deploy/k8s/teambudget-crd.yaml
```

### 2. Deploy RBAC

```bash
kubectl apply -f deploy/k8s/rbac.yaml
```

### 3. Deploy Controller and Agent

```bash
# Build Docker image
docker build -f deploy/docker/Dockerfile -t tco-configurator:latest .

# Deploy controller
kubectl apply -f deploy/k8s/controller-deployment.yaml

# Deploy agent
kubectl apply -f deploy/k8s/agent-daemonset.yaml
```

### 4. Create a Team Budget

```bash
kubectl apply -f deploy/k8s/example-teambudget.yaml
```

### 5. Verify Deployment

```bash
# Check pods
kubectl get pods

# Check TeamBudgets
kubectl get teambudgets

# View controller logs
kubectl logs -f deployment/tco-controller

# View agent logs
kubectl logs -f daemonset/tco-agent
```

## Deployment

### Production Deployment

#### 1. Build and Push Image

```bash
# Build
docker build -f deploy/docker/Dockerfile -t your-registry/tco-configurator:v1.0.0 .

# Push
docker push your-registry/tco-configurator:v1.0.0
```

#### 2. Update Manifests

Update image references in:
- `deploy/k8s/controller-deployment.yaml`
- `deploy/k8s/agent-daemonset.yaml`

#### 3. Deploy

```bash
kubectl apply -f deploy/k8s/
```

### Configuration

#### TeamBudget Resource

```yaml
apiVersion: tco.io/v1
kind: TeamBudget
metadata:
  name: backend-team
  namespace: default
spec:
  dailyLimit: 1000000      # 1MB daily (in bytes)
  monthlyLimit: 30000000   # 30MB monthly (in bytes)
status:
  currentUsage: 0          # Updated by agent
  lastUpdated: ""          # Timestamp of last update
```

#### Policy Thresholds

The policy engine uses these thresholds:

- **Allow:** `currentUsage + newBytes ≤ 80% of dailyLimit`
- **Throttle:** `80% < currentUsage + newBytes ≤ 100% of dailyLimit`
- **Drop:** `currentUsage + newBytes > dailyLimit`

### Monitoring

#### Prometheus Metrics

```
# Total actions by team and type
tco_policy_engine_actions_total{team="backend-team",action="allow"} 1234
tco_policy_engine_actions_total{team="backend-team",action="throttle"} 56
tco_policy_engine_actions_total{team="backend-team",action="drop"} 12
```

#### Dashboard

Access the dashboard:

```bash
# Run locally
go build -o dashboard-server ./dashboard
./dashboard-server

# Access at http://localhost:3000
```

## Development

### Project Structure

```
tco-configurator/
├── api/v1/                      # Kubernetes CRD definitions
│   ├── types.go                 # TeamBudget type definition
│   └── types_test.go            # CRD tests
├── cmd/                         # Binary entry points
│   ├── agent/main.go            # Agent service
│   ├── controller/main.go       # Controller service
│   ├── api/main.go              # REST API server
│   └── test-*/                  # Test utilities
├── dashboard/                   # Web UI
│   └── main.go                  # Dashboard server
├── deploy/                      # Deployment configs
│   ├── docker/Dockerfile        # Multi-stage build
│   └── k8s/                     # Kubernetes manifests
│       ├── teambudget-crd.yaml
│       ├── rbac.yaml
│       ├── controller-deployment.yaml
│       ├── agent-daemonset.yaml
│       └── example-teambudget.yaml
├── pkg/                         # Core packages
│   ├── agent/                   # Log processing agent
│   │   ├── agent.go
│   │   └── agent_test.go
│   ├── kubernetes/              # K8s client & controller
│   │   ├── client.go            # Dynamic client wrapper
│   │   ├── controller.go        # Watch & reconcile loop
│   │   └── controller_test.go
│   ├── policy/                  # Policy engine
│   │   ├── engine.go            # Stateless evaluation
│   │   ├── struct.go            # Budget types
│   │   ├── converter.go         # CRD to Budget conversion
│   │   └── *_test.go            # Comprehensive tests
│   ├── promtail/                # Promtail integration
│   │   └── plugin.go            # Log pipeline plugin
│   ├── notifier/                # Alert system
│   │   └── notifier.go          # Notification logic
│   └── metrics/                 # Prometheus metrics
│       └── metrics.go           # Metric definitions
└── readme.md                    # This file
```

### Building

```bash
# Build all binaries
go build ./cmd/agent
go build ./cmd/controller
go build ./cmd/api
go build -o dashboard-server ./dashboard

# Build Docker image
docker build -f deploy/docker/Dockerfile -t tco-configurator:latest .
```

### Testing

```bash
# Run all tests
go test ./... -v

# Run specific package tests
go test ./pkg/policy -v
go test ./pkg/agent -v
go test ./pkg/kubernetes -v

# Run with coverage
go test ./... -cover

# Run with race detection
go test ./... -race
```

### Test Coverage

```
Package                          Coverage
────────────────────────────────────────
pkg/policy                       85%
pkg/agent                        80%
pkg/kubernetes                   70%
pkg/promtail                     N/A (integration)
pkg/notifier                     N/A (simple)
────────────────────────────────────────
Overall                          78%
```

## API Reference

### TeamBudget CRD

```yaml
apiVersion: tco.io/v1
kind: TeamBudget
metadata:
  name: string              # Team identifier
  namespace: string         # K8s namespace
spec:
  dailyLimit: int64         # Daily limit in bytes
  monthlyLimit: int64       # Monthly limit in bytes
status:
  currentUsage: int64       # Current usage in bytes
  lastUpdated: string       # ISO 8601 timestamp
```

### Policy Engine API

```go
// Stateless evaluation function
func EvaluateUsageStateless(
    team string,           // Team identifier
    dailyLimit int64,      // Daily limit in bytes
    currentUsage int64,    // Current usage in bytes
    newBytes int64,        // New log size in bytes
) Action                   // Returns: Allow, Throttle, or Drop

// Action types
const (
    ActionAllow    Action = "allow"      // Pass log through
    ActionThrottle Action = "throttle"   // Sample log
    ActionDrop     Action = "drop"       // Discard log
)
```

### Agent API

```go
// Process logs and enforce policy
func (a *Agent) ProcessLogs(
    namespace string,      // K8s namespace (team identifier)
    logSize int64,         // Log entry size in bytes
) policy.Action            // Returns action taken
```

### REST API Endpoints

```
GET  /api/health                    # Health check
GET  /api/teams                     # List all teams
GET  /api/teams/{name}              # Get team details
GET  /metrics                       # Prometheus metrics
```

## Design Decisions

### Why Stateless Policy Engine?

**Decision:** Store state in K8s CRD status, not in memory

**Rationale:**
1. **Reliability:** Pod crashes don't lose state
2. **Scalability:** Multiple instances without coordination
3. **Simplicity:** No state synchronization needed
4. **Testability:** Pure functions are easy to test

**Trade-offs:**
-  Crash-safe
-  Horizontally scalable
-  Easy to test
-  Extra K8s API calls
-  Slight latency increase

### Why DaemonSet for Agent?

**Decision:** Deploy agent as DaemonSet, not Deployment

**Rationale:**
1. **Locality:** Process logs on the same node
2. **Scalability:** Automatically scales with nodes
3. **Reliability:** No single point of failure
4. **Performance:** Minimal network overhead

**Trade-offs:**
-  Low latency
-  Scales automatically
-  No SPOF
-  More pods (one per node)
-  Higher resource usage

### Why Kubernetes CRDs?

**Decision:** Use CRDs instead of ConfigMaps or external DB

**Rationale:**
1. **Native:** First-class K8s resources
2. **Declarative:** GitOps compatible
3. **RBAC:** Built-in access control
4. **Validation:** Schema validation
5. **Tooling:** kubectl, kustomize, helm support

**Trade-offs:**
-  K8s-native
-  GitOps ready
-  RBAC support
-  Requires CRD installation
-  K8s-specific

### Why Three-Tier Actions?

**Decision:** Allow → Throttle → Drop (not just Allow/Drop)

**Rationale:**
1. **Graceful:** Don't drop immediately
2. **Warning:** Throttle gives time to react
3. **Observability:** Maintain some logs even over budget
4. **UX:** Better user experience

**Trade-offs:**
-  Better UX
-  Maintains some observability
-  Time to react
-  More complex logic
-  Harder to predict exact costs

## Troubleshooting

### Pods Not Starting

```bash
# Check pod status
kubectl get pods

# View pod logs
kubectl logs <pod-name>

# Describe pod for events
kubectl describe pod <pod-name>

# Common issues:
# - Image pull errors: Check image name and registry
# - RBAC errors: Verify ServiceAccount and ClusterRole
# - CRD not installed: Apply teambudget-crd.yaml first
```

### Controller Not Reconciling

```bash
# Check controller logs
kubectl logs -f deployment/tco-controller

# Verify RBAC permissions
kubectl auth can-i get teambudgets --as=system:serviceaccount:default:tco-configurator

# Check if CRD is installed
kubectl get crd teambudgets.tco.io

# Verify TeamBudget exists
kubectl get teambudgets
```

### Agent Not Processing Logs

```bash
# Check agent logs
kubectl logs -f daemonset/tco-agent

# Verify agent can reach K8s API
kubectl exec -it <agent-pod> -- wget -O- http://kubernetes.default.svc/api

# Check RBAC permissions
kubectl auth can-i update teambudgets/status --as=system:serviceaccount:default:tco-configurator
```

### Metrics Not Appearing

```bash
# Check if metrics endpoint is accessible
kubectl port-forward deployment/tco-controller 8080:8080
curl http://localhost:8080/metrics

# Verify Prometheus is scraping
# Check Prometheus targets page
```


## License

**MIT License**

Copyright (c) 2025 Abhinav P Inamdar

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

---

**Built by Abhinav P Inamdar**

For questions or support, please open an issue on GitHub.

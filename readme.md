# TCO Configurator

Log volumes in cloud-native systems can grow uncontrollably, leading to unpredictable storage and processing costs.

TCO Configurator is a budgeting system that integrates with Promtail, Kubernetes, and Prometheus to:
- Track log volume per app/team
- Enforce budget limits via dynamic log filtering  
- Provide real-time visibility into cost vs. usage
- Help teams control costs without sacrificing observability

## Features

- **Log Volume Tracking** - Monitor usage per app/team/namespace
- **Budget Enforcement** - Drop/throttle logs when over budget
- **Prometheus Metrics** - Real-time cost dashboards
- **K8s Native** - CRDs for team-level budget policies

## Current Status

 **Policy Engine** - Core budget logic and evaluation  
 **Kubernetes CRDs** - TeamBudget resource definitions  
 **Testing** - Comprehensive test coverage  
 **Controller** - Watches TeamBudget resources (planned)  
 **Agent** - Modified Promtail for log tracking (planned)  
 **Dashboard** - Cost visualization UI (planned)  

## Quick Start

### 1. Install the CRD
```bash
kubectl apply -f deploy/k8s/teambudget-crd.yaml
```

### 2. Create a team budget
```bash
kubectl apply -f deploy/k8s/example-teambudget.yaml
```

### 3. Run tests
```bash
go test ./...
```

## Example Usage

Teams can define their log budgets using Kubernetes resources:

```yaml
apiVersion: tco.io/v1
kind: TeamBudget
metadata:
  name: backend-team
  namespace: default
spec:
  dailyLimit: 1000000    # 1MB daily
  monthlyLimit: 30000000 # 30MB monthly
```

The policy engine evaluates usage and returns actions:
- **Allow** - Under 80% of daily limit
- **Throttle** - 80-100% of daily limit  
- **Drop** - Over daily limit

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   TeamBudget    │───▶│  Policy Engine   │───▶│     Actions     │
│   (K8s CRD)     │    │  (Go Package)    │    │ (Allow/Throttle │
└─────────────────┘    └──────────────────┘    │    /Drop)       │
                                               └─────────────────┘
```

## Development

### Project Structure
```
├── api/v1/              # Kubernetes CRD definitions
├── pkg/policy/          # Core policy engine
├── deploy/k8s/          # Kubernetes manifests
├── cmd/                 # Main applications (planned)
└── test/                # Integration tests
```

### Building
```bash
make build
```

### Testing
```bash
make test
```

## License

**MIT License**

Copyright (c) 2025 Abhinav P Inamdar

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

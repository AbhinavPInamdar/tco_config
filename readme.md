Log volumes in cloud-native systems can grow uncontrollably, leading to unpredictable storage and processing costs.
TCO Configurator is a budgeting system that integrates with Promtail, Kubernetes, and Prometheus to:
    Track log volume per app/team.
    Enforce budget limits via dynamic log filtering.
    Provide real-time visibility into cost vs. usage.
    This helps teams control costs without sacrificing observability.


Features:
    Log Volume Tracking (per app/team/namespace)
    Budget Enforcement (drop/truncate logs when over budget)
    Prometheus Metrics for real-time cost dashboards
    K8s Native via CRDs for team-level budget policies


License:
MIT
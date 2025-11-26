package main

import (
	"context"
	"log"

	"tco-configurator/pkg/kubernetes"
	"tco-configurator/pkg/policy"
)

func main() {
	client, err := kubernetes.NewClient()
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	policyEngine := &policy.PolicyEngine{
		Teams: make(map[string]policy.Budget),
	}
	controller := kubernetes.Controller{
		KubeClient:   client,
		PolicyEngine: policyEngine,
	}
	ctx := context.Background()

	if err := controller.ReconcileTeambudget(ctx, "backend-team", "default"); err != nil {
		log.Fatalf("reconcile failed: %v", err)
	}
}

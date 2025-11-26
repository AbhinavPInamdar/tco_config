package main

import (
	"context"
	"log"
	"tco-configurator/pkg/kubernetes"
	"tco-configurator/pkg/policy"
)

func main() {
	kubeClient, err := kubernetes.NewClient()
	if err != nil {
		log.Fatalf("Failed to create K8s client: %v", err)
	}

	policyEngine := &policy.PolicyEngine{
		Teams: make(map[string]policy.Budget),
	}

	ctrl := &kubernetes.Controller{
		KubeClient:   kubeClient,
		PolicyEngine: policyEngine,
	}

	log.Println("Controller started")
	if err := ctrl.Start(context.Background(), ""); err != nil {
		log.Fatalf("controller.Start failed: %v", err)
	}
}

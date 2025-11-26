package main

import (
	"log"
	"tco-configurator/pkg/agent"
	"tco-configurator/pkg/kubernetes"
)

func main() {
	kubeClient, err := kubernetes.NewClient()
	if err != nil {
		log.Fatalf("Failed to create K8s client: %v", err)
	}

	_ = &agent.Agent{
		KubeClient: kubeClient,
	}

	log.Println("Agent started")
	select {}
}

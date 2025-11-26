package main

import (
	"fmt"
	"log"

	"tco-configurator/pkg/kubernetes"
)

func main() {
	client, err := kubernetes.NewClient()
	if err != nil {
		log.Fatalf("expected team budgets got %s", err)
	}
	teamb, err := client.ListTeamBudgets()
	if err != nil {
		log.Fatalf("failed to list teamBudgets: %v", err)
	}

	fmt.Printf("Found %d TeamBudgets\n", len(teamb))
}

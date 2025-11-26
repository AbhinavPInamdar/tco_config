package kubernetes

import (
	"context"
	"fmt"
	"tco-configurator/pkg/policy"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Controller struct {
	KubeClient   *Client
	PolicyEngine *policy.PolicyEngine
}

func (c *Controller) ReconcileTeambudget(ctx context.Context, teamBudgetName string, namespace string) error {
	TeamBudget, err := c.KubeClient.GetTeamBudget(teamBudgetName, namespace)
	if err != nil {
		return fmt.Errorf("expected team budget, got %v", err)
	}
	budget := policy.TeamBudgetToBudget(TeamBudget)
	c.PolicyEngine.AddTeamBudget(teamBudgetName, budget)
	result := c.PolicyEngine.EvaluateUsage(teamBudgetName, 0)
	fmt.Printf("Team %s action: %s\n", teamBudgetName, result)
	return nil
}


func (c *Controller) Start(ctx context.Context, namespace string) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		res := c.KubeClient.dynClient.Resource(c.KubeClient.gvr)

		var watchRes dynamic.ResourceInterface
		if namespace != "" {
			watchRes = res.Namespace(namespace)
		} else {
			watchRes = res
		}

		watcher, err := watchRes.Watch(ctx, metav1.ListOptions{})
		if err != nil {
			return fmt.Errorf("Failed to get teambudget watcher, %v", err)
		}

		for event := range watcher.ResultChan() {
			if err := ctx.Err(); err != nil {
				watcher.Stop()
				return err
			}

			if event.Type != watch.Added && event.Type != watch.Modified {
				continue
			}

			u, ok := event.Object.(*unstructured.Unstructured)
			if !ok {
				continue
			}

			name := u.GetName()
			ns := u.GetNamespace()

			if err := c.ReconcileTeambudget(ctx, name, ns); err != nil {
				fmt.Printf("reconcile failed for TeamBudget %s/%s: %v\n", ns, name, err)
			}
		}
	}
}
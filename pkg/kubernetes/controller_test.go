package kubernetes

import (
	"testing"
	"context"
	"tco-configurator/pkg/policy"
)

func TestController(t *testing.T) {
	policyEngine := &policy.PolicyEngine{
		Teams: make(map[string]policy.Budget),
	}
	controller := &Controller{
		PolicyEngine: policyEngine,
	}

	err := controller.reconcileTeambudget(context.Background(), "test-team")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
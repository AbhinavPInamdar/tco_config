package policy

import (
	"testing"
	v1 "tco-configurator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)


func TestConverter(t *testing.T) {
	teamBudget := &v1.TeamBudget {
		ObjectMeta: metav1.ObjectMeta{
			Name: "backend-team",
		},
		Spec: v1.TeamBudgetSpec{
			DailyLimit:1000000,
			MonthlyLimit:30000000,
		},
		Status: v1.TeamBudgetStatus{
			CurrentUsage:500000,
		},
	}

	budget := TeamBudgetToBudget(teamBudget)
	if budget.TeamName != "backend-team" {
		t.Errorf("expected TeamName 'backend-team', got %s", budget.TeamName)
	}
	if budget.DailyLimit != 1000000 {
		t.Errorf("Expected Dailylimit 1000000m, got %d", budget.DailyLimit)
	}
}
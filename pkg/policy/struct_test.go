package policy

import "testing"


func TestBudgetLimits(t *testing.T) {
	budget := Budget{
		TeamName: "team1",
		DailyLimit:1000,
		CurrentUsage:1200,
	}

	result := budget.IsOverDailyLimit()

	if result != true {
		t.Errorf("Expected true got %v", result)
	}
}


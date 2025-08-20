package policy

import "testing"

func TestEngine(t *testing.T) {
	engine := PolicyEngine {
		Teams: make(map[string]Budget),
	}


	engine.Teams["backend"] = Budget {
		TeamName: "backend",
		DailyLimit: 1000,
		CurrentUsage: 0,
	}
	result := engine.EvaluateUsage("backend", 500)
	if result != ActionAllow {
		t.Errorf("expected allow, got %v", result)
	}
}
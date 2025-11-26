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


func TestEvaluateUsageStateless(t *testing.T) {
	
	result := EvaluateUsageStateless("backend",1000,0, 500)
	if result != ActionAllow {
		t.Errorf("expected allow, got %v", result)
	}
	result = EvaluateUsageStateless("backend1",1000,0, 900)
	if result != ActionThrottle {
		t.Errorf("expected allow, got %v", result)
	}
	result = EvaluateUsageStateless("backend",1000,0, 1004)
	if result != ActionDrop {
		t.Errorf("expected allow, got %v", result)
	}
	
}
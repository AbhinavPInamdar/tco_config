package agent

import (
	"testing"
	"tco-configurator/api/v1"
	"tco-configurator/pkg/policy"

)






type MockKubeClient struct{
	name string
	namespace string
	dailyLimit int64
	currentUsage int64
}


func (m *MockKubeClient) GetTeamBudget(name, namespace string ) (*v1.TeamBudget, error) {
	return &v1.TeamBudget {
		Spec: v1.TeamBudgetSpec {
			DailyLimit: m.dailyLimit, 

		},

		Status: v1.TeamBudgetStatus{
			CurrentUsage: m.currentUsage,
		},
	}, nil

	
}


func (m *MockKubeClient) UpdateTeamBudgetStatus(name, namespace string, newUsage int64) error {
	return nil
}


func TestProcessLogs(t *testing.T) {
	mock := &MockKubeClient{
		dailyLimit: 1000,
		currentUsage: 500,
	}

	agent := &Agent {
		KubeClient: mock,
	}


	result := agent.ProcessLogs("default", 100)


	if result != policy.ActionAllow {
		t.Errorf("Expected allow, got %v", result)
	}
}


func TestProcessLogsThrottle(t *testing.T) {
	mock := &MockKubeClient{
		dailyLimit: 1000,
		currentUsage: 850,
	}

	agent := &Agent {
		KubeClient: mock,
	}


	result := agent.ProcessLogs("default", 100)


	if result != policy.ActionThrottle {
		t.Errorf("Expected throttle, got %v", result)
	}

}


func TestProcessLogsDrop(t *testing.T) {
	mock := &MockKubeClient{
		dailyLimit:1000,
		currentUsage: 950,

	}

	agent := &Agent {
		KubeClient: mock,
	}


	result := agent.ProcessLogs("default", 100)


	if result != policy.ActionDrop {
		t.Errorf("Expected drop, got %v", result)
	}
	
}





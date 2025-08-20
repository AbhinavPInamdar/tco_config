package v1
import (
	"testing"
	"io/ioutil"
	"sigs.k8s.io/yaml"
)


func TestTeamBudgetYAML(t *testing.T) {
	data,  err := ioutil.ReadFile("../../deploy/k8s/example-teambudget.yaml")
	if err != nil {
		t.Fatalf("failed to read YAML: %v", err)
	}
	var teamBudget TeamBudget

	err = yaml.Unmarshal(data, &teamBudget)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if teamBudget.Spec.DailyLimit != 1000000 {
		t.Errorf("expected dailylimit 1000000, got %d", teamBudget.Spec.DailyLimit)
	}
	if teamBudget.ObjectMeta.Name != "backend-team" {
		t.Errorf("expected name 'backend-team', got %s",teamBudget.ObjectMeta.Name )
	}

}
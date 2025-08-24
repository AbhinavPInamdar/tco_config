package policy

import "tco-configurator/pkg/metrics"

type PolicyEngine struct {
	Teams map[string]Budget
}

type Action string

const (
	ActionAllow    Action = "allow"
	ActionThrottle Action = "throttle"
	ActionAlert    Action = "alert"
	ActionDrop     Action = "drop"
)

func (pe *PolicyEngine) EvaluateUsage(team string, newUsage int64) Action {
	budget, ok := pe.Teams[team]
	if !ok {
		// This part was already correct.
		metrics.TeamBudgetActionsCounter.WithLabelValues(team, string(ActionDrop)).Inc()
		return ActionDrop
	}

	totalbytes := budget.CurrentUsage + newUsage
	var action Action


	switch {
	case totalbytes > budget.DailyLimit:
		action = ActionDrop
	case totalbytes > budget.DailyLimit*80/100:
		action = ActionThrottle
	default:
		action = ActionAllow
	}

	
	metrics.TeamBudgetActionsCounter.WithLabelValues(team, string(action)).Inc()

	return action
}

func (pe *PolicyEngine) AddTeamBudget(teamName string, budget Budget) {
	pe.Teams[teamName] = budget
}
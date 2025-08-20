package policy

type PolicyEngine struct {
	Teams map[string]Budget
}


type Action string


const (
	ActionAllow Action = "allow"
	ActionThrottle Action = "Throttle"
	ActionAlert Action = "alert"
	ActionDrop Action = "drop"
)


func (pe *PolicyEngine) EvaluateUsage(team string, newUsage int64) Action {
	budget ,  ok := pe.Teams[team]
	if !ok {
		return ActionDrop
	}
	totalbytes := budget.CurrentUsage+newUsage

	switch {
	case totalbytes > budget.DailyLimit:
		return ActionDrop
	case totalbytes > budget.DailyLimit*80/100: 
		return ActionThrottle
		
	default:
		return ActionAllow
	}

}

func (pe *PolicyEngine) AddTeamBudget(teamName string, budget Budget) {
	pe.Teams[teamName] = budget
}


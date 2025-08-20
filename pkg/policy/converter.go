package policy

import "tco-configurator/api/v1"

func TeamBudgetToBudget(tb *v1.TeamBudget) Budget {
	return Budget{
		TeamName: tb.ObjectMeta.Name,
		DailyLimit: tb.Spec.DailyLimit,
		MonthlyLimit: tb.Spec.MonthlyLimit,
		CurrentUsage: tb.Status.CurrentUsage,
	}
}
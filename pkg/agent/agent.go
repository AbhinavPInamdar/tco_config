package agent

import (
	"tco-configurator/pkg/kubernetes"
	"tco-configurator/pkg/policy"
)

type Agent struct {
	KubeClient   kubernetes.KubeClientInterface
	PolicyEngine *policy.PolicyEngine
}

func (a *Agent) ProcessLogs(namespace string, logSize int64) policy.Action {
	var action policy.Action
	budget, err := a.KubeClient.GetTeamBudget(namespace, namespace)
	if err != nil {
		return policy.ActionDrop
	}
	dailyLimit := budget.Spec.DailyLimit
	currUsage := budget.Status.CurrentUsage

	action = policy.EvaluateUsageStateless(namespace, dailyLimit, currUsage, logSize)
	newUsage := currUsage + logSize
	err = a.KubeClient.UpdateTeamBudgetStatus(namespace, namespace, newUsage)
	if err != nil {
		return policy.ActionDrop
	}

	return action
}

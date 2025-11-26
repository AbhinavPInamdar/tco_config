package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	TeamBudgetActionsCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tco_policy_engine_actions_total",
			Help: "The total number of actions taken by the policy engine, partitioned by team and action type.",
		},
		[]string{"team", "action"}, 
	)
)


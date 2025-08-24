package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// TeamBudgetActionsCounter is a counter for the actions taken by the policy engine
	TeamBudgetActionsCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tco_policy_engine_actions_total",
			Help: "The total number of actions taken by the policy engine, partitioned by team and action type.",
		},
		[]string{"team", "action"}, // Labels for partitioning the data
	)

	// Add more metrics here as you need them. For example, a gauge for current usage:
	// TeamCurrentUsage = promauto.NewGaugeVec(...)
)
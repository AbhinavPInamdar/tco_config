package notifier

import (
	"log"

	"tco-configurator/pkg/policy"
)

type Notifier struct {
}

func (n *Notifier) NotifyAction(team string, action policy.Action) {
	switch action {
	case policy.ActionDrop:
		log.Printf("ALERT: Team %s exceeded budget - dropping logs", team)
	case policy.ActionThrottle:
		log.Printf("WARNING: Team %s at 80%% of budget - throttling logs", team)
	}
}

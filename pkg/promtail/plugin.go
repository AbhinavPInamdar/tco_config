package promtail

import (
	"log"
	"time"

	"github.com/prometheus/common/model"
	"tco-configurator/pkg/agent"
	"tco-configurator/pkg/policy"
)

type PromtailPlugin struct {
	agent *agent.Agent
}

func NewPromtailPlugin(a *agent.Agent) *PromtailPlugin {
	return &PromtailPlugin{agent: a}
}

func (p *PromtailPlugin) Handle(labels model.LabelSet, t time.Time, entry string) error {
	namespace := string(labels["namespace"])

	if namespace == "" {
		namespace = "default"
	}

	logSize := int64(len(entry))

	action := p.agent.ProcessLogs(namespace, logSize)

	switch action {
	case policy.ActionAllow:
		log.Printf("Allow log from %s", namespace)
	case policy.ActionThrottle:
		log.Printf("Throttle log from %s", namespace)
	case policy.ActionDrop:
		log.Printf("Drop log from %s", namespace)
		return nil
	}
	return nil
}

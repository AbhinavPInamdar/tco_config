package v1

import "k8s.io/apimachinery/pkg/apis/meta/v1"



type TeamBudget struct{
	v1.TypeMeta `json:",inline"`
	v1.ObjectMeta `json:"metadata,omitempty"`
	Spec TeamBudgetSpec `json:"spec,omitempty"`
	Status TeamBudgetStatus `json:"status,omitempty"`
}

type TeamBudgetSpec struct {
	DailyLimit int64 `json:"dailyLimit"`
	MonthlyLimit int64 `json:"monthlyLimit"`

}

type TeamBudgetStatus struct {
	CurrentUsage int64 `json:"currentUsage"`
	LastUpdated string `json:"lastUpdated"`
}


type TeamBudgetList struct {
	v1.TypeMeta `json:",inline"`
	v1.ListMeta `json:"metadata,omitempty"`
	Items []TeamBudget `json:"items"`
}
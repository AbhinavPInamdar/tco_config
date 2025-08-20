package policy


type Budget struct {
	TeamName 		string
	DailyLimit  	int64
	MonthlyLimit 	int64
	CurrentUsage 	int64
}

func (b Budget) IsOverDailyLimit() bool {
	if b.CurrentUsage > b.DailyLimit {
		return true
	} else {
		return false
	}
}

func (b Budget) IsOverMonthlyLimit() bool {
	if b.CurrentUsage > b.MonthlyLimit {
		return true
	} else {
		return false
	}
}

func (b Budget) RemainingBudget() int64 {
	return b.DailyLimit-b.CurrentUsage
}

func (b Budget) PercentageUsage() float64 {
	return float64(b.CurrentUsage)/float64(b.DailyLimit)
}
 

package main

import (
	"encoding/json"
	"net/http"
	"tco-configurator/pkg/policy"
)

func teamsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    teams := []map[string]interface{}{
        {
            "name": "backend",
            "dailyLimit": 1000000,
            "currentUsage": 800000,
            "status": "allow",
        },
        {
            "name": "frontend", 
            "dailyLimit": 500000,
            "currentUsage": 600000,
            "status": "throttle",
        },
    }
    
    json.NewEncoder(w).Encode(teams)
}


func teamHandler(w http.ResponseWriter, r *http.Request){
	
}

func healthHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status":"ok"})
}


func main() {

	policyEngine:= &policy.PolicyEngine{
		Teams: make(map[string]policy.Budget),
	}
	policyEngine.AddTeamBudget("backend", policy.Budget{
    	TeamName:     "backend",
    	DailyLimit:   1000000,
    	MonthlyLimit: 30000000,
    	CurrentUsage: 800000,	
	})

	policyEngine.AddTeamBudget("frontend", policy.Budget{
    	TeamName:     "frontend", 
    	DailyLimit:   500000,
    	MonthlyLimit: 15000000,
    	CurrentUsage: 600000,
	})


	http.HandleFunc("/api/teams",teamsHandler)
	http.HandleFunc("/api/teams/{name}",teamHandler)
	http.HandleFunc("/api/health",healthHandler)
	http.ListenAndServe(":8080", nil)


}


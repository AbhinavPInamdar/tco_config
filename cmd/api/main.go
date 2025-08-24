package main

import (
	"encoding/json"
	"log"
	"net/http"
	"tco-configurator/pkg/policy"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func teamsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")


	teams := []map[string]interface{}{
		{
			"name":         "backend",
			"dailyLimit":   1000000,
			"currentUsage": 800000,
			"status":       "allow",
		},
		{
			"name":         "frontend",
			"dailyLimit":   500000,
			"currentUsage": 600000,
			"status":       "throttle",
		},
	}

	if err := json.NewEncoder(w).Encode(teams); err != nil {
		log.Printf("Error encoding teams data: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func teamHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamName := vars["name"]

	w.Header().Set("Content-Type", "application/json")
	log.Printf("Fetching data for team: %s", teamName)


	response := map[string]string{
		"teamName": teamName,
		"status":   "data_is_mocked",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding team data for %s: %v", teamName, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"status": "ok"}
	json.NewEncoder(w).Encode(response)
}

func main() {
	policyEngine := &policy.PolicyEngine{
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

	_ = policyEngine.EvaluateUsage("backend",0)
	_ = policyEngine.EvaluateUsage("frontend",100)

	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/health", healthHandler).Methods("GET")
	api.HandleFunc("/teams", teamsHandler).Methods("GET")
	api.HandleFunc("/teams/{name}", teamHandler).Methods("GET")
	r.Handle("/metrics", promhttp.Handler())

	log.Println("Starting server on port :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
package main

import (
	"log"
	"net/http"

	"tco-configurator/pkg/kubernetes"
)

func main() {
	client, err := kubernetes.NewClient()
	if err != nil {
		log.Fatalf("Failed to create K8s client: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		budgets, err := client.ListTeamBudgets()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<h1>TCO Configurator</h1>"))
		for _, b := range budgets {
			w.Write([]byte("<p>" + b.ObjectMeta.Name + "</p>"))
		}
	})

	log.Println("Dashboard running on :3000")
	http.ListenAndServe(":3000", nil)
}

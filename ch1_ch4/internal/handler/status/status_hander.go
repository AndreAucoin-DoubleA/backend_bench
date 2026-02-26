package status

import (
	"backend_bench/internal/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func StatusHandler(wikiRepo *model.WikiRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Example: Fetch some data from the wiki repository to demonstrate functionality
		wikiStats, err := wikiRepo.GetRecentChanges()
		if err != nil {
			fmt.Println("Error fetching recent changes:", err)
			http.Error(w, "Failed to fetch recent changes", http.StatusInternalServerError)
			return
		}

		resp := struct {
			TotalChanges  int            `json:"total_changes"`
			NumBots       int            `json:"num_bots"`
			NumNonBots    int            `json:"num_non_bots"`
			DistinctUsers int            `json:"distinct_users"`
			DistinctUrl   map[string]int `json:"distinct_url"`
		}{
			TotalChanges:  wikiStats.TotalChanges,
			NumBots:       wikiStats.NumBots,
			NumNonBots:    wikiStats.NumNonBots,
			DistinctUsers: wikiStats.DistinctUsers,
			DistinctUrl:   wikiStats.DistinctUrl,
		}

		data, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "Failed to encode stats", http.StatusInternalServerError)
			return
		}

		// Write headers and JSON body
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(data); err != nil {
			log.Println("Error writing response:", err)
		}
	}
}

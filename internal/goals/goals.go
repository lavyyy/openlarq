package goals

import (
	"encoding/json"
	"log"
	"maps"
	"net/http"
	"slices"

	"barking.dev/larq-api/internal/firebase"
)

type HydrationGoalEntry struct {
	Time          string  `json:"time"`
	Type          string  `json:"type"`
	VolumeInLiter float64 `json:"volumeInLiter"`
}

type HydrationGoalResponse struct {
	Entries []HydrationGoalEntry `json:"entries"`
}

func GetHydrationGoals(client *firebase.FirebaseClient) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()

		viewFrom := query.Get("viewFrom")
		index := query.Get("index")

		paramMap := map[string]interface{}{
			"viewFrom": viewFrom,
			"index":    index,
		}

		params := firebase.NewQueryParams(paramMap)

		res, err := client.GetUserHydrationGoals(params)
		if err != nil {
			log.Printf("Error getting hydration goals: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Format the response
		var formattedResponse HydrationGoalResponse
		formattedResponse.Entries = make([]HydrationGoalEntry, 0)

		// Iterate over each entry in the data, ensuring order is kept
		for _, key := range slices.Sorted(maps.Keys(res.Data.(map[string]interface{}))) {
			v := res.Data.(map[string]interface{})[key]
			if entryMap, ok := v.(map[string]interface{}); ok {
				var goalEntry HydrationGoalEntry
				if time, ok := entryMap["time"].(string); ok {
					goalEntry.Time = time
				}
				if typeStr, ok := entryMap["type"].(string); ok {
					goalEntry.Type = typeStr
				}
				if volume, ok := entryMap["volumeInLiter"].(float64); ok {
					goalEntry.VolumeInLiter = volume
				}
				formattedResponse.Entries = append(formattedResponse.Entries, goalEntry)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(formattedResponse); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

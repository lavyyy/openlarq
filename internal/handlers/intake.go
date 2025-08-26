package handlers

import (
	"encoding/json"
	"log"
	"maps"
	"net/http"
	"slices"

	"barking.dev/openlarq/internal/firebase"
)

type LiquidIntakeEntry struct {
	DateCreated   string  `json:"dateCreated"`
	Source        string  `json:"source"`
	Time          string  `json:"time"`
	Type          string  `json:"type"`
	VolumeInLiter float64 `json:"volumeInLiter"`
}

type LiquidIntakeResponse struct {
	Entries []LiquidIntakeEntry `json:"entries"`
}

func GetLiquidIntake(client *firebase.FirebaseClient) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()

		startTime := query.Get("startTime")
		endTime := query.Get("endTime")
		index := query.Get("index")

		paramMap := map[string]interface{}{
			"startTime": startTime,
			"endTime":   endTime,
			"index":     index,
		}

		params := firebase.NewQueryParams(paramMap)

		// use the new user-specific method
		res, err := client.GetUserLiquidIntake(params)
		if err != nil {
			log.Printf("Error getting liquid intake: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// format the response
		var formattedResponse LiquidIntakeResponse
		formattedResponse.Entries = make([]LiquidIntakeEntry, 0)

		// extract data from the response
		for _, key := range slices.Sorted(maps.Keys(res.Data.(map[string]interface{}))) {
			v := res.Data.(map[string]interface{})[key]

			if entryMap, ok := v.(map[string]interface{}); ok {
				var intakeEntry LiquidIntakeEntry
				if dateCreated, ok := entryMap["dateCreated"].(string); ok {
					intakeEntry.DateCreated = dateCreated
				}
				if source, ok := entryMap["source"].(string); ok {
					intakeEntry.Source = source
				}
				if time, ok := entryMap["time"].(string); ok {
					intakeEntry.Time = time
				}
				if typeStr, ok := entryMap["type"].(string); ok {
					intakeEntry.Type = typeStr
				}
				if volume, ok := entryMap["volumeInLiter"].(float64); ok {
					intakeEntry.VolumeInLiter = volume
				}
				formattedResponse.Entries = append(formattedResponse.Entries, intakeEntry)
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

package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"net/http"
	"slices"

	"barking.dev/openlarq/internal/cache"
	"barking.dev/openlarq/internal/firebase"
)

type HydrationGoalEntry struct {
	Time          string  `json:"time"`
	Type          string  `json:"type"`
	VolumeInLiter float64 `json:"volumeInLiter"`
}

type HydrationGoalResponse struct {
	Entries []HydrationGoalEntry `json:"entries"`
}

func GetHydrationGoals(client *firebase.FirebaseClient, cache *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()

		viewFrom := query.Get("viewFrom")
		index := query.Get("index")

		// create cache key
		cacheKey := fmt.Sprintf("hydration-goals:%s:%s", viewFrom, index)

		// try to get from cache first
		if cachedData, hit := cache.Get(cacheKey); hit {
			log.Printf("Cache hit for hydration goals, serving cached data")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(cachedData)
			return
		}

		// cache miss, fetch from Firebase
		response, err := fetchHydrationGoals(client, viewFrom, index)
		if err != nil {
			log.Printf("Error getting hydration goals: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// cache the response
		cache.Set(cacheKey, response)

		// return JSON response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func fetchHydrationGoals(client *firebase.FirebaseClient, viewFrom, index string) (*HydrationGoalResponse, error) {
	paramMap := map[string]interface{}{
		"viewFrom": viewFrom,
		"index":    index,
	}

	params := firebase.NewQueryParams(paramMap)

	res, err := client.GetUserHydrationGoals(params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch hydration goals: %w", err)
	}

	// format the response
	response := &HydrationGoalResponse{
		Entries: make([]HydrationGoalEntry, 0),
	}

	// iterate over each entry in the data, ensuring order is kept
	for _, key := range slices.Sorted(maps.Keys(res.Data.(map[string]interface{}))) {
		v := res.Data.(map[string]interface{})[key]
		if entryMap, ok := v.(map[string]interface{}); ok {
			entry := formatHydrationGoalEntry(entryMap)
			response.Entries = append(response.Entries, entry)
		}
	}

	return response, nil
}

func formatHydrationGoalEntry(entryMap map[string]interface{}) HydrationGoalEntry {
	var entry HydrationGoalEntry

	if time, ok := entryMap["time"].(string); ok {
		entry.Time = time
	}
	if typeStr, ok := entryMap["type"].(string); ok {
		entry.Type = typeStr
	}
	if volume, ok := entryMap["volumeInLiter"].(float64); ok {
		entry.VolumeInLiter = volume
	}

	return entry
}

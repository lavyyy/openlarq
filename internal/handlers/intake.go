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

func GetLiquidIntake(firebaseClient *firebase.FirebaseClient, cache *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()

		startTime := query.Get("startTime")
		endTime := query.Get("endTime")
		index := query.Get("index")

		// create cache key (v2 = exclude deleted entries)
		cacheKey := fmt.Sprintf("intake:v2:%s:%s:%s", startTime, endTime, index)

		// try to get from cache first
		if cachedData, hit := cache.Get(cacheKey); hit {
			log.Printf("Cache hit for liquid intake, serving cached data")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(cachedData)
			return
		}

		// cache miss, fetch from Firebase
		response, err := fetchLiquidIntake(firebaseClient, startTime, endTime, index)
		if err != nil {
			log.Printf("Error getting liquid intake: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// cache the response
		cache.Set(cacheKey, response)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func fetchLiquidIntake(firebaseClient *firebase.FirebaseClient, startTime, endTime, index string) (*LiquidIntakeResponse, error) {
	paramMap := map[string]interface{}{
		"startTime": startTime,
		"endTime":   endTime,
		"index":     index,
	}

	params := firebase.NewQueryParams(paramMap)

	res, err := firebaseClient.GetUserLiquidIntake(params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from Firebase: %w", err)
	}

	// format the response
	response := &LiquidIntakeResponse{
		Entries: make([]LiquidIntakeEntry, 0),
	}

	dataMap, _ := res.Data.(map[string]interface{})
	keys := slices.Sorted(maps.Keys(dataMap))

	for _, key := range keys {
		v := dataMap[key]
		entryMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		if isEntryDeleted(entryMap) {
			continue
		}
		entry := formatIntakeEntry(entryMap)
		response.Entries = append(response.Entries, entry)
	}

	return response, nil
}

// isEntryDeleted returns true if the intake entry is soft-deleted (updateState "deleted", dateDeleted set, or isDeleted true).
func isEntryDeleted(entryMap map[string]interface{}) bool {
	if u, ok := entryMap["updateState"].(string); ok && u == "deleted" {
		return true
	}
	if dateDeleted, ok := entryMap["dateDeleted"].(string); ok && dateDeleted != "" {
		return true
	}
	if isDeleted, ok := entryMap["isDeleted"].(bool); ok && isDeleted {
		return true
	}
	return false
}

func formatIntakeEntry(entryMap map[string]interface{}) LiquidIntakeEntry {
	var entry LiquidIntakeEntry

	if dateCreated, ok := entryMap["dateCreated"].(string); ok {
		entry.DateCreated = dateCreated
	}
	if source, ok := entryMap["source"].(string); ok {
		entry.Source = source
	}
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

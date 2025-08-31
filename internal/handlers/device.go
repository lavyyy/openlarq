package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"barking.dev/openlarq/internal/cache"
	"barking.dev/openlarq/internal/firebase"
)

type DeviceInfoResponse struct {
	Name                    string  `json:"name"`
	Color                   string  `json:"color"`
	SizeInMilliliter        float64 `json:"sizeInMilliliter"`
	PureVisPowerMode        string  `json:"pureVisPowerMode"`
	IsFilterTrackingEnabled bool    `json:"isFilterTrackingEnabled"`
}

func GetDeviceInfo(client *firebase.FirebaseClient, cache *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()
		deviceId := query.Get("deviceId")

		// create cache key
		cacheKey := fmt.Sprintf("device-info:%s", deviceId)

		// try to get from cache first
		if cachedData, hit := cache.Get(cacheKey); hit {
			log.Printf("Cache hit for device info, serving cached data")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(cachedData)
			return
		}

		// cache miss, fetch from Firebase
		response, err := fetchDeviceInfoFromFirebase(client, deviceId)
		if err != nil {
			log.Printf("Error getting device info: %v", err)
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

func fetchDeviceInfoFromFirebase(client *firebase.FirebaseClient, deviceId string) (*DeviceInfoResponse, error) {
	paramMap := map[string]interface{}{}
	params := firebase.NewQueryParams(paramMap)

	res, err := client.GetDeviceInfo(params, deviceId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch device info: %w", err)
	}

	response := &DeviceInfoResponse{}
	data := res.Data.(map[string]interface{})

	if name, ok := data["name"].(string); ok {
		response.Name = name
	}
	if color, ok := data["color"].(string); ok {
		response.Color = color
	}
	if size, ok := data["sizeInMilliliter"].(float64); ok {
		response.SizeInMilliliter = size
	}
	if powerMode, ok := data["pureVisPowerMode"].(string); ok {
		response.PureVisPowerMode = powerMode
	}
	if filterTracking, ok := data["isFilterTrackingEnabled"].(bool); ok {
		response.IsFilterTrackingEnabled = filterTracking
	}

	return response, nil
}

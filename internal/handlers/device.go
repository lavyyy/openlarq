package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"barking.dev/openlarq/internal/firebase"
)

type DeviceInfoResponse struct {
	Name                    string  `json:"name"`
	Color                   string  `json:"color"`
	SizeInMilliliter        float64 `json:"sizeInMilliliter"`
	PureVisPowerMode        string  `json:"pureVisPowerMode"`
	IsFilterTrackingEnabled bool    `json:"isFilterTrackingEnabled"`
}

func GetDeviceInfo(client *firebase.FirebaseClient) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		query := req.URL.Query()
		deviceId := query.Get("deviceId")

		paramMap := map[string]interface{}{}

		params := firebase.NewQueryParams(paramMap)

		res, err := client.GetDeviceInfo(params, deviceId)
		if err != nil {
			log.Printf("Error getting hydration goals: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		var formattedResponse DeviceInfoResponse

		data := res.Data.(map[string]interface{})

		if name, ok := data["name"].(string); ok {
			formattedResponse.Name = name
		}
		if color, ok := data["color"].(string); ok {
			formattedResponse.Color = color
		}
		if size, ok := data["sizeInMilliliter"].(float64); ok {
			formattedResponse.SizeInMilliliter = size
		}
		if powerMode, ok := data["pureVisPowerMode"].(string); ok {
			formattedResponse.PureVisPowerMode = powerMode
		}
		if filterTracking, ok := data["isFilterTrackingEnabled"].(bool); ok {
			formattedResponse.IsFilterTrackingEnabled = filterTracking
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(formattedResponse); err != nil {
			log.Printf("Error encoding response: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"barking.dev/openlarq/internal/firebase"
)

type UserInfoResponse struct {
	DisplayName string `json:"displayName"`
}

type CustomerDataResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func GetUserInfo(client *firebase.FirebaseClient) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		email := os.Getenv("LARQ_EMAIL")
		if email == "" {
			http.Error(w, "LARQ_EMAIL env variable not set", http.StatusInternalServerError)
			return
		}

		customerDataUrl := fmt.Sprintf("https://api.livelarq.com/api/v2/customer/%s", email)

		idToken := client.IdToken()
		if idToken == "" {
			http.Error(w, "No ID token available", http.StatusUnauthorized)
			return
		}

		customerReq, err := http.NewRequest("GET", customerDataUrl, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating request: %v", err), http.StatusInternalServerError)
			return
		}
		customerReq.Header.Set("Authorization", "Bearer "+idToken)

		// make GET request to the livelarq api
		httpClient := &http.Client{}
		resp, err := httpClient.Do(customerReq)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching customer data: %v", err), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading response body: %v", err), http.StatusInternalServerError)
			return
		}

		var customerData CustomerDataResponse
		if err := json.Unmarshal(body, &customerData); err != nil {
			http.Error(w, fmt.Sprintf("Error parsing customer data: %v", err), http.StatusInternalServerError)
			return
		}

		displayName := customerData.FirstName

		response := UserInfoResponse{
			DisplayName: displayName,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

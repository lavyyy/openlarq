package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"barking.dev/openlarq/internal/cache"
	"barking.dev/openlarq/internal/firebase"
)

type UserInfoResponse struct {
	DisplayName string `json:"displayName"`
}

type CustomerDataResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func GetUserInfo(client *firebase.FirebaseClient, cache *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// create cache key
		cacheKey := "user-info"

		// try to get from cache first
		if cachedData, hit := cache.Get(cacheKey); hit {
			log.Printf("Cache hit for user info, serving cached data")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(cachedData)
			return
		}

		// cache miss, fetch from API
		response, err := fetchUserInfo(client)
		if err != nil {
			log.Printf("Error getting user info: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// cache the response
		cache.Set(cacheKey, response)

		// return JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func fetchUserInfo(client *firebase.FirebaseClient) (*UserInfoResponse, error) {
	email := os.Getenv("LARQ_EMAIL")
	if email == "" {
		return nil, fmt.Errorf("LARQ_EMAIL env variable not set")
	}

	customerDataUrl := fmt.Sprintf("https://api.livelarq.com/api/v2/customer/%s", email)

	idToken := client.IdToken()
	if idToken == "" {
		return nil, fmt.Errorf("no ID token available")
	}

	customerReq, err := http.NewRequest("GET", customerDataUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	customerReq.Header.Set("Authorization", "Bearer "+idToken)

	// make GET request to the livelarq api
	httpClient := &http.Client{}
	resp, err := httpClient.Do(customerReq)
	if err != nil {
		return nil, fmt.Errorf("error fetching customer data: %w", err)
	}
	defer resp.Body.Close()

	// read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var customerData CustomerDataResponse
	if err := json.Unmarshal(body, &customerData); err != nil {
		return nil, fmt.Errorf("error parsing customer data: %w", err)
	}

	displayName := customerData.FirstName

	response := &UserInfoResponse{
		DisplayName: displayName,
	}

	return response, nil
}

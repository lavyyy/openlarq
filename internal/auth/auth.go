package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type AuthResponse struct {
	Kind         string `json:"kind"`
	LocalId      string `json:"localId"`
	Email        string `json:"email"`
	DisplayName  string `json:"displayName"`
	IdToken      string `json:"idToken"`
	Registered   bool   `json:"registered"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
}

func Authenticate() (string, error) {
	// create request body
	reqBody := map[string]interface{}{
		"clientType":        "CLIENT_TYPE_IOS",
		"email":             os.Getenv("LARQ_EMAIL"),
		"password":          os.Getenv("LARQ_PASSWORD"),
		"returnSecureToken": true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request body: %v", err)
	}

	authUrl := "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyPassword?key=AIzaSyDHKFqlodeFaDa2sQqxxBAaW-A3JF5lXRU"

	res, err := http.Post(authUrl, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	responseData := &AuthResponse{}
	err = json.NewDecoder(res.Body).Decode(responseData)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return "", err
	}

	return responseData.IdToken, nil
}

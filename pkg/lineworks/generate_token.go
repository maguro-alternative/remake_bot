package lineworks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const baseAuthURL = "https://auth.worksmobile.com/oauth2/v2.0"

type lineWorksTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    string    `json:"expires_in"`
	Scope 	  string `json:"scope"`
}

func (l LineWorksInfo) GetAccessToken(ctx context.Context, scope string) (*lineWorksTokenResponse, error) {
	// Generate JWT
	jwt, err := generateJWT(l.lineWorksClientID, l.lineWorksServiceAccount, l.lineWorksPrivateKey)
	if err != nil {
		return nil, err
	}

	// Prepare request
	tokenURL := fmt.Sprintf("%s/token", baseAuthURL)
	formData := url.Values{
		"assertion":     {jwt},
		"grant_type":    {"urn%3Aietf%3Aparams%3Aoauth%3Agrant-type%3Ajwt-bearer"},
		"client_id":     {l.lineWorksClientID},
		"client_secret": {l.lineWorksClientSecret},
		"scope":         {scope},
	}

	// Make POST request
	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := l.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode response body
	var body *lineWorksTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetAccessToken: %v", resp.Status)
	}

	return body, nil
}

func (l LineWorksInfo) RefreshAccessToken(ctx context.Context, refreshToken string) (*lineWorksTokenResponse, error) {
	formData := url.Values{
		"refresh_token": {refreshToken},
		"grant_type":    {"refresh_token"},
		"client_id":     {l.lineWorksClientID},
		"client_secret": {l.lineWorksClientSecret},
	}
	tokenURL := fmt.Sprintf("%s/token", baseAuthURL)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := l.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode response body
	var body *lineWorksTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	return body, nil
}

func generateJWT(clientID, serviceAccount, privateKey string) (string, error) {
	currentTime := time.Now().Unix()
	claims := jwt.MapClaims{
		"iss": clientID,
		"sub": serviceAccount,
		"iat": currentTime,
		"exp": currentTime + 3600, // 1 hour
	}

	// Parse the RSA private key
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return "", err
	}

	// Create the token and sign it
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	jws, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return jws, nil
}

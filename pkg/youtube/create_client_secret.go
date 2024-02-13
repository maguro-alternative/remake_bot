package youtube

import (
	"encoding/json"
)

type clientSecret struct {
	Installed struct {
		ClientID                string   `json:"client_id"`
		ProjectID               string   `json:"project_id"`
		AuthUri                 string   `json:"auth_uri"`
		TokenUri                string   `json:"token_uri"`
		AuthProviderX509CertUrl string   `json:"auth_provider_x509_cert_url"`
		ClientSecret            string   `json:"client_secret"`
		RedirectUris            []string `json:"redirect_uris"`
	} `json:"installed"`
}

func createClientSecret(youtubeClientID, youtubeProjectID, youtubeClientSecret string) ([]byte, error) {
	clientData := clientSecret{
		Installed: struct {
			ClientID                string   `json:"client_id"`
			ProjectID               string   `json:"project_id"`
			AuthUri                 string   `json:"auth_uri"`
			TokenUri                string   `json:"token_uri"`
			AuthProviderX509CertUrl string   `json:"auth_provider_x509_cert_url"`
			ClientSecret            string   `json:"client_secret"`
			RedirectUris            []string `json:"redirect_uris"`
		}{
			ClientID:                youtubeClientID,
			ProjectID:               youtubeProjectID,
			AuthUri:                 "https://accounts.google.com/o/oauth2/auth",
			TokenUri:                "https://oauth2.googleapis.com/token",
			AuthProviderX509CertUrl: "https://www.googleapis.com/oauth2/v1/certs",
			ClientSecret:            youtubeClientSecret,
			RedirectUris:            []string{"http://localhost"},
		},
	}
	return json.Marshal(clientData)
}
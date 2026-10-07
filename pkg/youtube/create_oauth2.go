package youtube

import (
	"encoding/json"

	"github.com/samber/mo"
)

type oAuth2Credentials struct {
	AccessToken   string            `json:"access_token"`
	ClientID      string            `json:"client_id"`
	ClientSecret  string            `json:"client_secret"`
	RefreshToken  string            `json:"refresh_token"`
	TokenExpiry   string            `json:"token_expiry"`
	TokenURI      string            `json:"token_uri"`
	UserAgent     mo.Option[string] `json:"user_agent"`
	RevokeURI     string            `json:"revoke_uri"`
	IDToken       mo.Option[string] `json:"id_token"`
	IDTokenJWT    mo.Option[string] `json:"id_token_jwt"`
	TokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		Scope       string `json:"scope"`
		TokenType   string `json:"token_type"`
	} `json:"token_response"`
	Scopes       []string `json:"scopes"`
	TokenInfoURI string   `json:"token_info_uri"`
	Invalid      bool     `json:"invalid"`
	Class        string   `json:"_class"`
	Module       string   `json:"_module"`
}

func createOAuth2(accessToken, clientID, clientSecret, refreshToken, tokenExpiry string) ([]byte, error) {
	oauth2Data := oAuth2Credentials{
		AccessToken:  accessToken,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RefreshToken: refreshToken,
		TokenExpiry:  tokenExpiry,
		TokenURI:     "https://oauth2.googleapis.com/token",
		UserAgent:    mo.None[string](),
		RevokeURI:    "https://oauth2.googleapis.com/revoke",
		IDToken:      mo.None[string](),
		IDTokenJWT:   mo.None[string](),
		TokenResponse: struct {
			AccessToken string `json:"access_token"`
			ExpiresIn   int    `json:"expires_in"`
			Scope       string `json:"scope"`
			TokenType   string `json:"token_type"`
		}{
			AccessToken: accessToken,
			ExpiresIn:   3599,
			Scope:       "https://www.googleapis.com/auth/youtube.upload",
			TokenType:   "Bearer",
		},
		Scopes:       []string{"https://www.googleapis.com/auth/youtube.upload"},
		TokenInfoURI: "https://oauth2.googleapis.com/tokeninfo",
		Invalid:      false,
		Class:        "OAuth2Credentials",
		Module:       "oauth2client.client",
	}
	return json.Marshal(oauth2Data)
}

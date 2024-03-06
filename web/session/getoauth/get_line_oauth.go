package getoauth

import (
	"errors"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/gorilla/sessions"
)

type LineOAuthSession struct {
	Token oauth2.Token `json:"token"`
	User  LineUser     `json:"user"`
}

type LineToken struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type LineUser struct {
	Iss      string   `json:"iss"`
	Sub      string   `json:"sub"`
	Aud      string   `json:"aud"`
	Exp      int      `json:"exp"`
	Iat      int      `json:"iat"`
	AuthTime int      `json:"auth_time"`
	Nonce    int      `json:"nonce"`
	Amr      []string `json:"amr"`
	Name     string   `json:"name"`
	Picture  string   `json:"picture"`
	Email    string   `json:"email"`
}

func GetLineOAuth(store *sessions.CookieStore, r *http.Request, sessionSecret string) (*LineUser, error) {
	session, err := store.Get(r, sessionSecret)
	if err != nil {
		return nil, err
	}
	// セッションに保存されているlineuserを取得
	lineUser, ok := session.Values["line_user"].(*LineUser)
	if !ok {
		return nil, errors.New("session not found")
	}
	return lineUser, nil
}
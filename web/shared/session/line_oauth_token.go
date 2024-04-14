package session

import (
	"errors"
)

var lineOAuthTokenKey sessionKey = "line_oauth_token"

func (s *sessionStore) SetLineOAuthToken(token string) {
	s.session.Values[lineOAuthTokenKey] = token
}

func (s *sessionStore) GetLineOAuthToken() (string, error) {
	token, ok := s.session.Values[lineOAuthTokenKey].(string)
	if !ok {
		return "", errors.New("line oauth token not found")
	}
	return token, nil
}

func (s *sessionStore) CleanupLineOAuthToken() {
	delete(s.session.Values, lineOAuthTokenKey)
}

package session

import (
	"errors"
)

var csrfTokenKey sessionKey = "csrf_token"

func (s *sessionStore) SetCSRFToken(token string) {
	s.session.Values[csrfTokenKey] = token
}

func (s *sessionStore) GetCSRFToken() (string, error) {
	token, ok := s.session.Values[csrfTokenKey].(string)
	if !ok || token == "" {
		return "", errors.New("csrf token not found")
	}
	return token, nil
}

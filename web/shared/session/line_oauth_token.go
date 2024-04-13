package session

import (
	"context"
)

var lineOAuthTokenKey = "line_oauth_token"

func (s *SessionStore) SetLineOAuthToken(ctx context.Context, token string) {
	s.session.Values[lineOAuthTokenKey] = token
}

func (s *SessionStore) GetLineOAuthToken() (string, bool) {
	token, ok := s.session.Values[lineOAuthTokenKey].(string)
	return token, ok
}

func (s *SessionStore) CleanupLineOAuthToken(ctx context.Context) {
	delete(s.session.Values, lineOAuthTokenKey)
}

package session

import (
	"context"
)

var lineOAuthTokenKey sessionKey = "line_oauth_token"

func (s *sessionStore) SetLineOAuthToken(ctx context.Context, token string) {
	s.session.Values[lineOAuthTokenKey] = token
}

func (s *sessionStore) GetLineOAuthToken() (string, bool) {
	token, ok := s.session.Values[lineOAuthTokenKey].(string)
	return token, ok
}

func (s *sessionStore) CleanupLineOAuthToken(ctx context.Context) {
	delete(s.session.Values, lineOAuthTokenKey)
}

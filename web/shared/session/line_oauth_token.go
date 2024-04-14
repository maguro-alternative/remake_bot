package session

var lineOAuthTokenKey sessionKey = "line_oauth_token"

func (s *sessionStore) SetLineOAuthToken(token string) {
	s.session.Values[lineOAuthTokenKey] = token
}

func (s *sessionStore) GetLineOAuthToken() (string, bool) {
	token, ok := s.session.Values[lineOAuthTokenKey].(string)
	return token, ok
}

func (s *sessionStore) CleanupLineOAuthToken() {
	delete(s.session.Values, lineOAuthTokenKey)
}

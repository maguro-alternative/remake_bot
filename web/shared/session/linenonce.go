package session

var lineNonceKey sessionKey = "line_nonce"

func (s *sessionStore) SetLineNonce(nonce string) {
	s.session.Values[lineNonceKey] = nonce
}

func (s *sessionStore) GetLineNonce() (string, bool) {
	state, ok := s.session.Values[lineNonceKey].(string)
	return state, ok
}

func (s *sessionStore) CleanupLineNonce() {
	delete(s.session.Values, lineNonceKey)
}

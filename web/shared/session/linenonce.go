package session

import (
	"context"
)

var lineNonceKey sessionKey = "line_nonce"

func (s *SessionStore) SetLineNonce(ctx context.Context, nonce string) {
	s.session.Values[lineNonceKey] = nonce
}

func (s *SessionStore) GetLineNonce() (string, bool) {
	state, ok := s.session.Values[lineNonceKey].(string)
	return state, ok
}

func (s *SessionStore) CleanupLineNonce(ctx context.Context) {
	delete(s.session.Values, lineNonceKey)
}

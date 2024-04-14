package session

import (
	"context"
)

var lineNonceKey sessionKey = "line_nonce"

func (s *sessionStore) SetLineNonce(ctx context.Context, nonce string) {
	s.session.Values[lineNonceKey] = nonce
}

func (s *sessionStore) GetLineNonce() (string, bool) {
	state, ok := s.session.Values[lineNonceKey].(string)
	return state, ok
}

func (s *sessionStore) CleanupLineNonce(ctx context.Context) {
	delete(s.session.Values, lineNonceKey)
}

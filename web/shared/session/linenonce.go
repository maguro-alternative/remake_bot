package session

import (
	"errors"
)

var lineNonceKey sessionKey = "line_nonce"

func (s *sessionStore) SetLineNonce(nonce string) {
	s.session.Values[lineNonceKey] = nonce
}

func (s *sessionStore) GetLineNonce() (string, error) {
	state, ok := s.session.Values[lineNonceKey].(string)
	if !ok {
		return "", errors.New("line nonce not found")
	}
	return state, nil
}

func (s *sessionStore) CleanupLineNonce() {
	delete(s.session.Values, lineNonceKey)
}

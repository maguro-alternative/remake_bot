package session

import (
	"errors"
)

var lineStateKey sessionKey = "line_state"

func (s *sessionStore) SetLineState(state string) {
	s.session.Values[lineStateKey] = state
}

func (s *sessionStore) GetLineState() (string, error) {
	state, ok := s.session.Values[lineStateKey].(string)
	if !ok {
		return "", errors.New("line state not found")
	}
	return state, nil
}

func (s *sessionStore) CleanupLineState() {
	s.session.Values[lineStateKey] = ""
	delete(s.session.Values, lineStateKey)
}

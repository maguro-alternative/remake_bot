package session

import (
	"context"
)

var lineStateKey sessionKey = "line_state"

func (s *SessionStore) SetLineState(ctx context.Context, state string) {
	s.session.Values[lineStateKey] = state
}

func (s *SessionStore) GetLineState() (string, bool) {
	state, ok := s.session.Values[lineStateKey].(string)
	return state, ok
}

func (s *SessionStore) CleanupLineState(ctx context.Context) {
	delete(s.session.Values, lineStateKey)
}

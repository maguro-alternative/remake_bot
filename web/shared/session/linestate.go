package session

import (
	"context"
)

var lineStateKey sessionKey = "line_state"

func (s *sessionStore) SetLineState(ctx context.Context, state string) {
	s.session.Values[lineStateKey] = state
}

func (s *sessionStore) GetLineState() (string, bool) {
	state, ok := s.session.Values[lineStateKey].(string)
	return state, ok
}

func (s *sessionStore) CleanupLineState(ctx context.Context) {
	delete(s.session.Values, lineStateKey)
}

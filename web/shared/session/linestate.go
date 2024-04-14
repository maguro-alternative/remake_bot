package session



var lineStateKey sessionKey = "line_state"

func (s *sessionStore) SetLineState(state string) {
	s.session.Values[lineStateKey] = state
}

func (s *sessionStore) GetLineState() (string, bool) {
	state, ok := s.session.Values[lineStateKey].(string)
	return state, ok
}

func (s *sessionStore) CleanupLineState() {
	delete(s.session.Values, lineStateKey)
}

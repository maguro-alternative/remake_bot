package session

var discordStateKey sessionKey = "discord_state"

func (s *sessionStore) SetDiscordState(state string) {
	s.session.Values[discordStateKey] = state
}

func (s *sessionStore) GetDiscordState() (string, bool) {
	state, ok := s.session.Values[discordStateKey].(string)
	return state, ok
}

func (s *sessionStore) CleanupDiscordState() {
	delete(s.session.Values, discordStateKey)
}

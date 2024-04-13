package session

import (
	"context"
)

var discordStateKey = "discord_state"

func (s *SessionStore) SetDiscordState(ctx context.Context, state string) {
	s.session.Values[discordStateKey] = state
}

func (s *SessionStore) GetDiscordState() (string, bool) {
	state, ok := s.session.Values[discordStateKey].(string)
	return state, ok
}

func (s *SessionStore) CleanupDiscordState(ctx context.Context) {
	delete(s.session.Values, discordStateKey)
}

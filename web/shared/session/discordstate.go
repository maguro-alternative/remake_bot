package session

import (
	"context"
)

var discordStateKey sessionKey = "discord_state"

func (s *sessionStore) SetDiscordState(ctx context.Context, state string) {
	s.session.Values[discordStateKey] = state
}

func (s *sessionStore) GetDiscordState() (string, bool) {
	state, ok := s.session.Values[discordStateKey].(string)
	return state, ok
}

func (s *sessionStore) CleanupDiscordState(ctx context.Context) {
	delete(s.session.Values, discordStateKey)
}

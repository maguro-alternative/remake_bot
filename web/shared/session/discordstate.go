package session

import (
	"errors"
)

var discordStateKey sessionKey = "discord_state"

func (s *sessionStore) SetDiscordState(state string) {
	s.session.Values[discordStateKey] = state
}

func (s *sessionStore) GetDiscordState() (string, error) {
	state, ok := s.session.Values[discordStateKey].(string)
	if !ok {
		return "", errors.New("discord state not found")
	}
	return state, nil
}

func (s *sessionStore) CleanupDiscordState() {
	s.session.Values[discordStateKey] = ""
	delete(s.session.Values, discordStateKey)
}

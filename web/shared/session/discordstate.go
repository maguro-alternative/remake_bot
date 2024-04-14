package session

import (
	"errors"

	"github.com/gorilla/sessions"
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
	delete(s.session.Values, discordStateKey)
}

func SetDiscordState(session *sessions.Session, state string) {
    session.Values[discordStateKey] = state
}

func GetDiscordState(session *sessions.Session) (string, error) {
    state, ok := session.Values[discordStateKey].(string)
    if !ok {
        return "", errors.New("discord state not found in session")
    }
    return state, nil
}

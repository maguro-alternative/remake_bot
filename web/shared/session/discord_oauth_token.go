package session

import (
	"errors"
)

var discordOAuthTokenKey sessionKey = "discord_oauth_token"

func (s *sessionStore) SetDiscordOAuthToken(token string) {
	s.session.Values[discordOAuthTokenKey] = token
}

func (s *sessionStore) GetDiscordOAuthToken() (string, error) {
	token, ok := s.session.Values[discordOAuthTokenKey].(string)
	if !ok {
		return "", errors.New("discord oauth token not found")
	}
	return token, nil
}

func (s *sessionStore) CleanupDiscordOAuthToken() {
	delete(s.session.Values, discordOAuthTokenKey)
}

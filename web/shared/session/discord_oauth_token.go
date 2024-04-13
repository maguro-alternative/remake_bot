package session

import (
	"context"
)

var discordOAuthTokenKey = "discord_oauth_token"

func (s *SessionStore) SetDiscordOAuthToken(ctx context.Context, token string) {
	s.session.Values[discordOAuthTokenKey] = token
}

func (s *SessionStore) GetDiscordOAuthToken() (string, bool) {
	token, ok := s.session.Values[discordOAuthTokenKey].(string)
	return token, ok
}

func (s *SessionStore) CleanupDiscordOAuthToken(ctx context.Context) {
	delete(s.session.Values, discordOAuthTokenKey)
}

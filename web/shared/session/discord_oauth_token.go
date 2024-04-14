package session

import (
	"context"
)

var discordOAuthTokenKey sessionKey = "discord_oauth_token"

func (s *sessionStore) SetDiscordOAuthToken(ctx context.Context, token string) {
	s.session.Values[discordOAuthTokenKey] = token
}

func (s *sessionStore) GetDiscordOAuthToken() (string, bool) {
	token, ok := s.session.Values[discordOAuthTokenKey].(string)
	return token, ok
}

func (s *sessionStore) CleanupDiscordOAuthToken(ctx context.Context) {
	delete(s.session.Values, discordOAuthTokenKey)
}

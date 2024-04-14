package session

var discordOAuthTokenKey sessionKey = "discord_oauth_token"

func (s *sessionStore) SetDiscordOAuthToken(token string) {
	s.session.Values[discordOAuthTokenKey] = token
}

func (s *sessionStore) GetDiscordOAuthToken() (string, bool) {
	token, ok := s.session.Values[discordOAuthTokenKey].(string)
	return token, ok
}

func (s *sessionStore) CleanupDiscordOAuthToken() {
	delete(s.session.Values, discordOAuthTokenKey)
}

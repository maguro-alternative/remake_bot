package session

var guildIdKey sessionKey = "guild_id"

func (s *sessionStore) SetGuildID(guildId string) {
	s.session.Values[guildId] = guildId
}

func (s *sessionStore) GetGuildID() (string, bool) {
	guildId, ok := s.session.Values[guildIdKey].(string)
	return guildId, ok
}

func (s *sessionStore) CleanupGuildID() {
	delete(s.session.Values, guildIdKey)
}

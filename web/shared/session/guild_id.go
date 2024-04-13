package session

import (
	"context"
)

var guildIdKey sessionKey = "guild_id"

func (s *SessionStore) SetGuildID(ctx context.Context, guildId string) {
	s.session.Values[guildId] = guildId
}

func (s *SessionStore) GetGuildID() (string, bool) {
	guildId, ok := s.session.Values[guildIdKey].(string)
	return guildId, ok
}

func (s *SessionStore) CleanupGuildID(ctx context.Context) {
	delete(s.session.Values, guildIdKey)
}

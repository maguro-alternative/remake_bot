package session

import (
	"errors"
)

var guildIdKey sessionKey = "guild_id"

func (s *sessionStore) SetGuildID(guildId string) {
	s.session.Values[guildId] = guildId
}

func (s *sessionStore) GetGuildID() (string, error) {
	guildId, ok := s.session.Values[guildIdKey].(string)
	if !ok {
		return "", errors.New("guild id not found")
	}
	return guildId, nil
}

func (s *sessionStore) CleanupGuildID() {
	delete(s.session.Values, guildIdKey)
}

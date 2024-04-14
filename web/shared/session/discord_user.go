package session

import (
	"github.com/maguro-alternative/remake_bot/web/shared/model"
)

var discordUserKey sessionKey = "discord_user"

func (s *sessionStore) SetDiscordUser(user *model.DiscordUser) {
	s.session.Values[discordUserKey] = user
}

func (s *sessionStore) GetDiscordUser() (*model.DiscordUser, bool) {
	user, ok := s.session.Values[discordUserKey].(*model.DiscordUser)
	return user, ok
}

func (s *sessionStore) CleanupDiscordUser() {
	delete(s.session.Values, discordUserKey)
}

package session

import (
	"errors"

	"github.com/maguro-alternative/remake_bot/web/shared/model"
)

var discordUserKey sessionKey = "discord_user"

func (s *sessionStore) SetDiscordUser(user *model.DiscordUser) {
	s.session.Values[discordUserKey] = user
}

func (s *sessionStore) GetDiscordUser() (*model.DiscordUser, error) {
	user, ok := s.session.Values[discordUserKey].(*model.DiscordUser)
	if !ok {
		return nil, errors.New("discord user not found")
	}
	return user, nil
}

func (s *sessionStore) CleanupDiscordUser() {
	delete(s.session.Values, discordUserKey)
}

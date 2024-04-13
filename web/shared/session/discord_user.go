package session

import (
	"context"

	"github.com/maguro-alternative/remake_bot/web/shared/session/model"
)

var discordUserKey sessionKey = "discord_user"

func (s *SessionStore) SetDiscordUser(ctx context.Context, user *model.DiscordUser) {
	s.session.Values[discordUserKey] = user
}

func (s *SessionStore) GetDiscordUser() (*model.DiscordUser, bool) {
	user, ok := s.session.Values[discordUserKey].(*model.DiscordUser)
	return user, ok
}

func (s *SessionStore) CleanupDiscordUser(ctx context.Context) {
	delete(s.session.Values, discordUserKey)
}

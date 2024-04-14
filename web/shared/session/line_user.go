package session

import (
	"context"

	"github.com/maguro-alternative/remake_bot/web/shared/model"
)

var lineUserKey sessionKey = "line_user"

func (s *sessionStore) SetLineUser(ctx context.Context, user model.LineIdTokenUser) {
	s.session.Values[lineUserKey] = user
}

func (s *sessionStore) GetLineUser(ctx context.Context) (model.LineIdTokenUser, bool) {
	user, ok := s.session.Values[lineUserKey].(model.LineIdTokenUser)
	return user, ok
}

func (s *sessionStore) ClearLineUser(ctx context.Context) {
	delete(s.session.Values, lineUserKey)
}

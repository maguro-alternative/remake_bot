package session

import (
	"github.com/maguro-alternative/remake_bot/web/shared/model"
)

var lineUserKey sessionKey = "line_user"

func (s *sessionStore) SetLineUser(user model.LineIdTokenUser) {
	s.session.Values[lineUserKey] = user
}

func (s *sessionStore) GetLineUser() (model.LineIdTokenUser, bool) {
	user, ok := s.session.Values[lineUserKey].(model.LineIdTokenUser)
	return user, ok
}

func (s *sessionStore) ClearLineUser() {
	delete(s.session.Values, lineUserKey)
}

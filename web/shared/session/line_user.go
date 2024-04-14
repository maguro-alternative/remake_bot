package session

import (
	"errors"

	"github.com/maguro-alternative/remake_bot/web/shared/model"
)

var lineUserKey sessionKey = "line_user"

func (s *sessionStore) SetLineUser(user *model.LineIdTokenUser) {
	s.session.Values[lineUserKey] = user
}

func (s *sessionStore) GetLineUser() (*model.LineIdTokenUser, error) {
	user, ok := s.session.Values[lineUserKey].(*model.LineIdTokenUser)
	if !ok {
		return nil, errors.New("line user not found")
	}
	return user, nil
}

func (s *sessionStore) CleanupLineUser() {
	s.session.Values[lineUserKey] = &model.LineIdTokenUser{}
	delete(s.session.Values, lineUserKey)
}

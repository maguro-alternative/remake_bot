package session

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type sessionKey string

type sessionStore struct {
	session *sessions.Session
}

func NewSessionStore(
	r *http.Request,
	store *sessions.CookieStore,
	secret string,
) *sessionStore {
	session, err := store.Get(r, secret)
	if err != nil {
		return nil
	}
	return &sessionStore{
		session: session,
	}
}

func (s *sessionStore) SessionSave(r *http.Request, w http.ResponseWriter) error {
	return s.session.Save(r, w)
}

func (s *sessionStore) GetSession() *sessions.Session {
	return s.session
}

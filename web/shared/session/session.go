package session

import (
	"net/http"

	"github.com/gorilla/sessions"
)

type sessionKey string

type SessionStore struct {
	session *sessions.Session
	store   *sessions.CookieStore
}

func NewSessionStore(
	r *http.Request,
	secret string,
) *SessionStore {
	store := sessions.NewCookieStore([]byte(secret))
	session, err := store.Get(r, secret)
	if err != nil {
		return nil
	}
	return &SessionStore{
		session: session,
		store:   store,
	}
}

func (s *SessionStore) SessionSave(w http.ResponseWriter) error {
	return s.session.Save(nil, w)
}

func (s *SessionStore) CookieSave(
	r *http.Request,
	w http.ResponseWriter,
) error {
	return s.store.Save(r, w, s.session)
}

package session

import (
	"encoding/gob"
	"net/http"

	"github.com/gorilla/sessions"
)

type sessionKey string

func init() {
	gob.Register(sessionKey(""))
}

type sessionStore struct {
	session *sessions.Session
}

func NewSessionStore(
	r *http.Request,
	store *sessions.CookieStore,
	secret string,
) (*sessionStore, error) {
	session, err := store.Get(r, secret)
	if err != nil {
		return nil, err
	}
	return &sessionStore{
		session: session,
	}, nil
}

func (s *sessionStore) SessionSave(r *http.Request, w http.ResponseWriter) error {
	return s.session.Save(r, w)
}

func (s *sessionStore) StoreSave(
	r *http.Request,
	w http.ResponseWriter,
	store *sessions.CookieStore,
) error {
	return store.Save(r, w, s.session)
}

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

// NewSessionStore retrieves or creates a session using the given session name.
// IMPORTANT: The sessionName parameter should be a non-secret identifier for the cookie,
// NOT the session secret key. The secret key should only be passed to sessions.NewCookieStore.
func NewSessionStore(
	r *http.Request,
	store *sessions.CookieStore,
	sessionName string,
) (*sessionStore, error) {
	session, err := store.Get(r, sessionName)
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

package session

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
)

type Session struct {
	Id string
}

const sessionIdNumBytes = 64

// Create a new Session with a random Id.
func NewSession() (*Session, error) {
	sess := &Session{}

	b := make([]byte, sessionIdNumBytes)
	_, err := rand.Read(b)
	if err != nil {
		return sess, fmt.Errorf("failed to generate a new session ID: %w", err)
	}

	sess.Id = hex.EncodeToString(b)

	return sess, nil
}

const cookieName = "session-id"

// GetOrCreate returns the session from the request, or generates a new session when none
// is present.
func GetOrCreate(r *http.Request) (*Session, bool, error) {
	id, err := r.Cookie(cookieName)
	if err == nil {
		return &Session{Id: id.Value}, false, nil
	}
	if !errors.Is(err, http.ErrNoCookie) {
		return nil, false, fmt.Errorf("failed to get the request session: %w", err)
	}

	sess, err := NewSession()
	if err != nil {
		return sess, false, fmt.Errorf("failed to create a new session: %w", err)
	}
	return sess, true, nil
}

func (s Session) Cookie() *http.Cookie {
	// NOTE(sean): switch to gorillatoolkit.org/pkg/securecookie
	return &http.Cookie{
		Name:     cookieName,
		Value:    s.Id,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
	}
}

package session

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
)

type Session struct {
	Id string
}

const sessionIdNumBytes = 64

// Create a new Session with a random Id.
func NewSession() (*Session, error) {
	s := &Session{Id: ""}

	b := make([]byte, sessionIdNumBytes)
	_, err := rand.Read(b)
	if err != nil {
		return s, fmt.Errorf("Failed to generate random session ID: %s", err)
	}

	s.Id = hex.EncodeToString(b)

	return s, nil
}

func (s Session) Cookie() *http.Cookie {
	// NOTE(sean): switch to gorillatoolkit.org/pkg/securecookie
	return &http.Cookie{
		Name:     "session-id",
		Value:    s.Id,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
	}
}

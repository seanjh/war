package session

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/seanjh/war/internal/appcontext"
)

type Session struct {
	ID string
}

const sessionIdNumBytes = 16

const cookieName = "session-id"

func generateSessionID() (string, error) {
	bytes := make([]byte, sessionIdNumBytes)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GetOrCreate returns the session from the request, or generates a new session when none
// is present.
func GetOrCreate(w http.ResponseWriter, r *http.Request) (*Session, error) {
	id, err := r.Cookie(cookieName)
	if err == nil {
		return &Session{ID: id.Value}, nil
	}
	if !errors.Is(err, http.ErrNoCookie) {
		return nil, fmt.Errorf("failed to get the request session: %w", err)
	}

	sessionId, err := generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to create new session ID: %w", err)
	}

	ctx := appcontext.GetAppContext(r)
	dbSess, err := ctx.DBWriter.Query.CreateSession(r.Context(), sessionId)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new session: %w", err)
	}

	sess := &Session{ID: dbSess.ID}
	http.SetCookie(w, sess.cookie())
	return sess, nil
}

func (s Session) cookie() *http.Cookie {
	// NOTE(sean): maybe switch to gorillatoolkit.org/pkg/securecookie
	return &http.Cookie{
		Name:     cookieName,
		Value:    s.ID,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
	}
}

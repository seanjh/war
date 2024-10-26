package session

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/seanjh/war/internal/appcontext"
	"github.com/seanjh/war/internal/db"
)

type Session struct {
	ID string
}

const sessionIdNumBytes = 64

const cookieName = "session-id"

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

	ctx := appcontext.GetAppContext(r)
	q := db.New(ctx.WriteDB)
	dbSess, err := q.CreateSession(r.Context())
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

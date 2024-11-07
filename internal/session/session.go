package session

import (
	"context"
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

const cookieName = "session-id"

func extractSessionID(r *http.Request) (string, error) {
	id, err := r.Cookie(cookieName)
	if err == nil {
		return id.Value, nil
	}
	if !errors.Is(err, http.ErrNoCookie) {
		return "", fmt.Errorf("failed to get the request session: %w", err)
	}

	return "", nil
}

func validateSessionId(sessionID string, r *http.Request) (*Session, error) {
	ctx := appcontext.GetAppContext(r)
	row, err := ctx.DBReader.Query.GetSession(r.Context(), sessionID)
	if err != nil {
		return nil, fmt.Errorf("session ID '%s' not recognized: %w", sessionID, err)
	}
	return &Session{ID: row.ID}, nil
}

const sessionIdNumBytes = 16

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
func OpenNewSession(w http.ResponseWriter, r *http.Request) (*Session, error) {
	ctx := appcontext.GetAppContext(r)
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to create new session ID: %w", err)
	}
	ctx.Logger.Info("Generated new random session ID",
		"sessionID", sessionID)

	dbSess, err := ctx.DBWriter.Query.CreateSession(r.Context(), sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new session: %w", err)
	}
	ctx.Logger.Info("Created new session",
		"sessionId", dbSess.ID, "created", dbSess.Created)

	http.SetCookie(w, cookie(dbSess.ID))
	return &Session{ID: dbSess.ID}, nil
}

func cookie(sessionID string) *http.Cookie {
	// NOTE(sean): maybe switch to gorillatoolkit.org/pkg/securecookie
	return &http.Cookie{
		Name:     cookieName,
		Value:    sessionID,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Secure:   true,
	}
}

const sessionIDKey string = "sessionid"

func WithSessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := extractSessionID(r)
		if err != nil {
			return
		}
		sess, err := validateSessionId(sessionID, r)
		if err != nil {
			return
		}
		ctx := context.WithValue(r.Context(), sessionIDKey, sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetSession(r *http.Request) *Session {
	sess := &Session{}
	sessionID, ok := r.Context().Value(sessionIDKey).(string)
	if ok {
		sess.ID = sessionID
	}
	return sess
}

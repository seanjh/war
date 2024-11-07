package game

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/seanjh/war/internal/appcontext"
	"github.com/seanjh/war/internal/db"
	"github.com/seanjh/war/internal/session"
)

type Player struct {
	Deck Deck
	Role GameRole
}

type GameRole int64

const (
	Unknown GameRole = iota
	Host
	Guest
)

func (r GameRole) String() string {
	switch r {
	case Host:
		return "host"
	case Guest:
		return "guest"
	}
	return "unknown"
}

func ConvertGameRole(val int64) GameRole {
	if val > int64(Guest) {
		return Unknown
	}
	return GameRole(val)
}

type Battle struct {
	Battle map[string]Card
	War    map[string][]Card
}

type Game struct {
	ID      int
	Player1 *Player
	Player2 *Player
	Battle  *Battle
}

// OpenNewGame returns a new Game with 2 Players with equal cuts of a new Deck.
func OpenNewGame(r *http.Request, sessionID string) (*Game, error) {
	ctx := appcontext.GetAppContext(r)

	tx, err := ctx.DBWriter.DB.Begin()
	defer tx.Rollback()
	if err != nil {
		return nil, fmt.Errorf("Failed to create new game: %w", err)
	}

	gameRow, err := ctx.DBWriter.Query.WithTx(tx).CreateGame(r.Context())
	if err != nil {
		return nil, fmt.Errorf("Failed to create new game: %w", err)
	}
	ctx.Logger.Info("Created new game row",
		"gameID", gameRow.ID,
		"gameCode", gameRow.Code)

	deck := NewDeck()
	deck.Shuffle(NewRiffleShuffler())
	d1, d2 := deck.Cut()

	err = ctx.DBWriter.Query.WithTx(tx).CreateHostGameSession(r.Context(), db.CreateHostGameSessionParams{
		GameID:    gameRow.ID,
		GameID_2:  gameRow.ID,
		Deck:      d1.String(),
		Deck_2:    d2.String(),
		SessionID: sql.NullString{String: sessionID, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create new host game session: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("Failed to commit new host game session: %w", err)
	}

	game := &Game{
		ID:      int(gameRow.ID),
		Player1: &Player{Deck: d1, Role: Host},
		Player2: &Player{Deck: d2, Role: Guest},
		Battle:  &Battle{},
	}
	return game, nil
}

func CreateFlip() http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles(
		filepath.Join("templates", "game.html"),
		filepath.Join("templates", "player.html"),
		filepath.Join("templates", "battleground.html"),
		filepath.Join("templates", "warzone.html"),
	))
	return func(w http.ResponseWriter, r *http.Request) {
		game, ok := mustHaveGame(w, "1")
		if !ok {
			http.Error(w, fmt.Sprintf("Cannot locate game: %s", "1"), http.StatusBadRequest)
			return
		}
		tmpl.ExecuteTemplate(w, "game", game)
	}
}

func Load(id string) (*Game, error) {
	game, ok := games[id]
	if !ok {
		return nil, fmt.Errorf("Game not found: %s", id)
	}
	return game, nil
}

// LoadGame returns a pre-existing Game when recognized.
func LoadGame(id string, r *http.Request) (*Game, error) {
	gameID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert gameID '%s' to int: %w", id, err)
	}

	game := &Game{ID: gameID}

	ctx := appcontext.GetAppContext(r)
	sess := session.GetSession(r)
	rows, err := ctx.DBReader.Query.GetGameSessions(r.Context(), int64(gameID))
	if err != nil {
		return nil, fmt.Errorf("Failed to load gameID '%d' from database: %w", gameID, err)
	}
	for _, row := range rows {
		role := ConvertGameRole(row.Role)
		deck := strings.Split(row.Deck, ",")
		switch role {
		case Host:
			game.Player1 = &Player{Role: Host, Deck: Deck(deck)}
		case Guest:
			game.Player2 = &Player{Role: Guest, Deck: Deck(deck)}
		default:
			ctx.Logger.Error("Unsupported player role",
				"sessionID", sess.ID,
				"row", row)
		}
	}
	return game, nil
}

func loadGameTemplates() *template.Template {
	return template.Must(template.ParseFiles(
		filepath.Join("templates", "layout.html"),
		filepath.Join("templates", "game.html"),
		filepath.Join("templates", "player.html"),
		filepath.Join("templates", "battleground.html"),
		filepath.Join("templates", "warzone.html"),
	))
}

func CreateAndRenderGame() http.HandlerFunc {
	tmpl := loadGameTemplates()
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := appcontext.GetAppContext(r)

		// sessionID, err := session.OpenNewSession(w, r)
		s := session.GetSession(r)
		if s.ID == "" {
			newSession, err := session.OpenNewSession(w, r)
			s = newSession
			if err != nil {
				ctx.Logger.Info("Failed to open new session",
					"err", err,
				)
				http.Error(w, "Failed to create new session", http.StatusInternalServerError)
				return
			}
		}

		game, err := OpenNewGame(r, s.ID)
		if err != nil {
			ctx.Logger.Error("Failed to create new game", "err", err)
			http.Error(w, "Failed to create new game", http.StatusInternalServerError)
			return
		}
		ctx.Logger.Info("Created new game and host game session",
			"gameID", game.ID,
		)
		w.Header().Add("hx-push-url", fmt.Sprintf("/game/%d", game.ID))
		if err := tmpl.ExecuteTemplate(w, "layout", game); err != nil {
			ctx.Logger.Error("Failed to render game template",
				"err", err,
				"gameID", game.ID,
				"sessionID", s.ID)
			http.Error(w, "Failed to render game template", http.StatusInternalServerError)
			return
		}
	}
}

func RenderGame() http.HandlerFunc {
	tmpl := loadGameTemplates()
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		game, ok := mustHaveGame(w, id)
		if !ok {
			http.Error(w, fmt.Sprintf("Cannot locate game: %s", id), http.StatusBadRequest)
			return
		}
		tmpl.ExecuteTemplate(w, "layout", game)
	}
}

func RenderHome() http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles(
		filepath.Join("templates", "layout.html"),
		filepath.Join("templates", "home.html"),
	))
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "layout", nil)
	}
}

// TODO(sean) remove this global state
// Temporary global game instances
var games map[string]*Game

func SetupRoutes(mux *http.ServeMux) *http.ServeMux {
	// TODO(sean) remove this global state
	games = make(map[string]*Game)

	mux.HandleFunc("GET /", RenderHome())
	mux.HandleFunc("POST /game", CreateAndRenderGame())
	mux.HandleFunc("GET /game/{id}", RenderGame())
	mux.HandleFunc("POST /flip", CreateFlip())
	session.WithSessionMiddleware(mux)

	return mux
}

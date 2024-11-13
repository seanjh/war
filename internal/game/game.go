package game

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

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
		return nil, fmt.Errorf("failed to create new game: %w", err)
	}

	gameRow, err := ctx.DBWriter.Query.WithTx(tx).CreateGame(r.Context())
	if err != nil {
		return nil, fmt.Errorf("failed to create new game: %w", err)
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
		return nil, fmt.Errorf("failed to create new host game session: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit new host game session: %w", err)
	}

	game := &Game{
		ID:      int(gameRow.ID),
		Player1: &Player{Deck: d1, Role: Host},
		Player2: &Player{Deck: d2, Role: Guest},
		Battle:  &Battle{},
	}
	return game, nil
}

// LoadGame returns a pre-existing Game when recognized.
func LoadGame(rawGameID string, r *http.Request) (*Game, error) {
	gameID, err := strconv.Atoi(rawGameID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert gameID '%s' to int: %w", rawGameID, err)
	}

	game := &Game{ID: gameID}

	ctx := appcontext.GetAppContext(r)
	rows, err := ctx.DBReader.Query.GetGameSessions(r.Context(), int64(gameID))
	if err != nil {
		return nil, fmt.Errorf("failed to load gameID '%d' from database: %w", gameID, err)
	}
	for _, row := range rows {
		role := ConvertGameRole(row.Role)
		deck, err := ConvertDeck(row.Deck)
		if err != nil {
			return nil, fmt.Errorf("failed to load deck for game ID %d: %w", gameID, err)
		}
		switch role {
		case Host:
			game.Player1 = &Player{Role: Host, Deck: deck}
		case Guest:
			game.Player2 = &Player{Role: Guest, Deck: deck}
		default:
			return nil, fmt.Errorf("unsupported player role %s", role)
		}
	}
	return game, nil
}

type PlayerContext struct {
	GameID     int
	Player     *Player
	BattleCard *Card
}

type GameContext struct {
	Player1 PlayerContext
	Player2 PlayerContext
}

func NewGameContext(game *Game) GameContext {
	return GameContext{
		Player1: PlayerContext{
			GameID:     game.ID,
			Player:     game.Player1,
			BattleCard: &Card{Suit: SuitHeart, Value: 2}},
		Player2: PlayerContext{
			GameID:     game.ID,
			Player:     game.Player2,
			BattleCard: &Card{Suit: SuitSpade, Value: 4}},
	}
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

func mustLoadGame(w http.ResponseWriter, r *http.Request) (*Game, error) {
	ctx := appcontext.GetAppContext(r)

	gameID := r.PathValue("id")
	ctx.UpdateLogger(ctx.Logger.With("gameID", gameID))

	game, err := LoadGame(gameID, r)
	if err != nil {
		ctx.Logger.Error("Failed to load game",
			"err", err,
			"gameID", gameID)
		http.Error(w, "cannot locate game", http.StatusBadRequest)
		return nil, fmt.Errorf("cannot load game with ID %s: %w", gameID, err)
	}
	return game, nil
}

func CreateAndRenderGame() http.HandlerFunc {
	tmpl := loadGameTemplates()
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := appcontext.GetAppContext(r)

		sess, err := session.LoadRequestSession(r)
		if err != nil {
			newSession, err := session.OpenNewSession(w, r)
			sess = newSession
			if err != nil {
				ctx.Logger.Error("Failed to open new session",
					"err", err)
				http.Error(w, "Failed to create game", http.StatusInternalServerError)
				return
			}
		}

		game, err := OpenNewGame(r, sess.ID)
		if err != nil {
			ctx.Logger.Error("Failed to open new game",
				"err", err)
			http.Error(w, "Failed to create game", http.StatusInternalServerError)
			return
		}
		w.Header().Add("hx-push-url", fmt.Sprintf("/game/%d", game.ID))

		ctx.Logger.Info("Created new game and host game session",
			"gameID", game.ID)

		data := GameContext{
			Player1: PlayerContext{GameID: game.ID, Player: game.Player1},
			Player2: PlayerContext{GameID: game.ID, Player: game.Player2},
		}
		if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
			ctx.Logger.Error("Failed to render game template",
				"err", err,
				"gameID", game.ID,
				"sessionID", sess.ID)
			http.Error(w, "Failed to render game", http.StatusInternalServerError)
			return
		}
	}
}

func RenderGame() http.HandlerFunc {
	tmpl := loadGameTemplates()
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := appcontext.GetAppContext(r)

		game, err := mustLoadGame(w, r)
		if err != nil {
			return
		}

		data := NewGameContext(game)
		err = tmpl.ExecuteTemplate(w, "layout", data)
		if err != nil {
			ctx.Logger.Error("ExecuteTemplate failed",
				"err", err)
			http.Error(w, "failed to load game", http.StatusInternalServerError)
			return
		}
	}
}

func CreateFlip() http.HandlerFunc {
	tmpl := loadGameTemplates()
	return func(w http.ResponseWriter, r *http.Request) {
		game, err := mustLoadGame(w, r)
		if err != nil {
			return
		}

		ctx := appcontext.GetAppContext(r)
		data := NewGameContext(game)
		err = tmpl.ExecuteTemplate(w, "layout", data)
		if err != nil {
			ctx.Logger.Error("ExecuteTemplate failed",
				"err", err,
				"gameID", game.ID)
			http.Error(w, "failed to load game", http.StatusInternalServerError)
			return
		}
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

func SetupRoutes(mux *http.ServeMux) *http.ServeMux {
	mux.Handle("GET /", http.HandlerFunc(RenderHome()))
	mux.Handle("POST /game", CreateAndRenderGame())
	mux.Handle("GET /game/{id}", session.RequireSession(RenderGame()))
	mux.Handle("POST /game/{id}/flip", session.RequireSession(CreateFlip()))
	return mux
}

package game

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/seanjh/war/internal/appcontext"
	"github.com/seanjh/war/internal/db"
	"github.com/seanjh/war/internal/session"
)

type Player struct {
	CardsInHand         Deck
	CardsWon            Deck
	CardsInBattle       Deck
	HiddenCardsInBattle []Deck
	Role                GameRole
	Session             session.Session
}

type War struct {
	Hand       Deck
	Battling   Deck
	Supporting []Deck
}

func (w *War) flip() error {
	if len(w.Hand) == 0 {
		return fmt.Errorf("no cards available to flip")
	}

	// When no other cards are battling, this is the opening of a simple 1 card play
	if len(w.Battling) == 0 {
		card, hand := w.Hand[0], w.Hand[1:]
		w.Hand = hand
		w.Battling = append(w.Battling, card)
		return nil
	}

	// take: 1 - non war
	// war take: 1 - B=1 (hand len 1), 2 - S=0 B=1 (hand len 2), 3 S=0 S=1 B=2 (hand len 3), 4 - S=0 S=1 S=2 B=3 (hand len >=4)

	// take: 0 (hand len 0), 1 (hand len 1), 2 (hand len 2), or 3 (hand len >=3)

	// When there are already cards battling, the next play is a war, where up to a
	// maximum of 3 cards are played in support (when possible), and 1 card is added to
	// the battling deck.
	take := int(math.Min(float64(len(w.Hand)), 4.0)) // 1-4
	supporting := make(Deck, take-1)
	for i := 0; i < take-1; i++ {
		supporting[i] = w.Hand[i]
		w.Hand = w.Hand[1:]
	}
	// card := w.Hand[0]
	// w.Hand = w.Hand[1:]
	// supporting, battling, hand := w.Hand[], w.Hand[0], w.Hand[take:]

	if take == 1 {
		card, hand := w.Hand[0], w.Hand[1:]
		w.Hand = hand
		w.Battling = append(w.Battling, card)
		return nil
	}

	cards, hand := w.Hand[:take], w.Hand[take:]
	w.Supporting = append(w.Supporting, cards)
	w.Hand = hand

	return nil
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

type Game struct {
	ID      int
	Player1 *Player
	Player2 *Player
}

// OpenNewGame returns a new Game with 2 Players with equal cuts of a new Deck.
func OpenNewGame(r *http.Request, sessionID string) (*Game, error) {
	ctx := appcontext.GetAppContext(r)

	tx, err := ctx.DBWriter.DB.BeginTx(r.Context(), &sql.TxOptions{})
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
		GameID:        gameRow.ID,
		GameID_2:      gameRow.ID,
		CardsInHand:   d1.String(),
		CardsInHand_2: d2.String(),
		SessionID:     sql.NullString{String: sessionID, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create new host game session: %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit new host game session: %w", err)
	}

	game := &Game{
		ID: int(gameRow.ID),
		Player1: &Player{
			CardsInHand: d1,
			Role:        Host,
			Session:     session.Session{ID: sessionID},
		},
		Player2: &Player{
			CardsInHand: d2,
			Role:        Guest,
			Session:     session.Session{},
		},
	}
	return game, nil
}

func NewGameFromGameSessionRows(gameID int, rows []db.GetGameSessionsRow) (*Game, error) {
	game := &Game{ID: gameID}
	for _, row := range rows {
		role := ConvertGameRole(row.Role)
		deck, err := ConvertDeck(row.CardsInHand)
		if err != nil {
			return nil, fmt.Errorf("failed to load deck for game ID %d: %w", gameID, err)
		}
		switch role {
		case Host:
			game.Player1 = &Player{
				Role:        Host,
				CardsInHand: deck,
				Session:     session.Session{ID: row.SessionID},
			}
		case Guest:
			game.Player2 = &Player{
				Role:        Guest,
				CardsInHand: deck,
				Session:     session.Session{ID: row.SessionID},
			}
		default:
			return nil, fmt.Errorf("unsupported player role %s", role)
		}
	}
	return game, nil
}

// LoadGame returns a pre-existing Game when recognized.
func LoadGame(rawGameID string, r *http.Request) (*Game, error) {
	gameID, err := strconv.Atoi(rawGameID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert gameID '%s' to int: %w", rawGameID, err)
	}

	ctx := appcontext.GetAppContext(r)
	rows, err := ctx.DBReader.Query.GetGameSessions(r.Context(), int64(gameID))
	if err != nil {
		return nil, fmt.Errorf("failed to load gameID '%d' from database: %w", gameID, err)
	}
	return NewGameFromGameSessionRows(gameID, rows)
}

func (game *Game) createFlip(r *http.Request, gameID int, session session.Session) error {
	ctx := appcontext.GetAppContext(r)

	tx, err := ctx.DBReader.DB.BeginTx(r.Context(), &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	row, err := ctx.DBWriter.Query.WithTx(tx).GetGameSession(
		r.Context(),
		db.GetGameSessionParams{GameID: 1, SessionID: sql.NullString{String: session.ID, Valid: true}})
	if err != nil {
		return fmt.Errorf("failed to read game session: %w", err)
	}
	hand, err := ConvertDeck(row.CardsInHand)
	if err != nil {
		return fmt.Errorf("failed to convert cards in hand to deck: %w", err)
	}
	if len(hand) < 1 {
		return errors.New("no cards available to flip")
	}
	card, hand := hand[len(hand)-1], hand[:len(hand)-1]

	inBattle, err := ConvertDeck(row.CardsInBattle)
	if err != nil {
		return fmt.Errorf("failed to convert cards in battle to deck: %w", err)
	}
	inBattle = append(inBattle, card)

	err = ctx.DBReader.Query.WithTx(tx).FlipCard(r.Context(), db.FlipCardParams{
		ID:            row.ID,
		CardsInHand:   hand.String(),
		CardsInBattle: inBattle.String()})
	if err != nil {
	}

	return nil
}

type PlayerContext struct {
	GameID     int
	Player     *Player
	BattleCard *Card
	TotalCards int
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
			BattleCard: &Card{Suit: SuitHeart, Value: 2},
			TotalCards: len(game.Player1.CardsInHand) + len(game.Player1.CardsWon),
		},
		Player2: PlayerContext{
			GameID:     game.ID,
			Player:     game.Player2,
			BattleCard: &Card{Suit: SuitSpade, Value: 4},
			TotalCards: len(game.Player2.CardsInHand) + len(game.Player2.CardsWon),
		},
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

		data := NewGameContext(game)
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
		sess := session.GetSession(r)

		err = game.createFlip(r, game.ID, sess)
		if err != nil {
			ctx.Logger.Error("Failed to perform card flip",
				"err", err)
			http.Error(w, "cannot flip", http.StatusInternalServerError)
			return
		}

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

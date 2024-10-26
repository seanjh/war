package game

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/seanjh/war/internal/appcontext"
	"github.com/seanjh/war/internal/session"
)

type Player struct {
	Deck Deck
	Name string
	Id   string
}

// NewPlayer returns a new player.
func NewPlayer(d Deck, id string, name string) *Player {
	p := Player{Deck: d, Id: id, Name: name}
	return &p
}

type Battle struct {
	Battle map[string]Card
	War    map[string][]Card
}

type Game struct {
	Id      string
	Player1 *Player
	Player2 *Player
	Battle  *Battle
}

// NewGame returns a new Game with 2 Players with equal cuts of a new Deck.
func NewGame() *Game {
	deck := NewDeck()
	deck.Shuffle(NewRiffleShuffler())
	d1, d2 := deck.Cut()
	return &Game{
		Id:      "1",
		Player1: NewPlayer(d1, "1", "Player One"),
		Player2: NewPlayer(d2, "2", "Player Two"),
		Battle:  &Battle{},
	}
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

func mustHaveGame(w http.ResponseWriter, id string) (*Game, bool) {
	game, err := Load(id)
	if err != nil {
		return &Game{}, false
	}
	return game, true
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
		game := NewGame()
		games[game.Id] = game

		ctx := appcontext.GetAppContext(r)
		sess, err := session.GetOrCreate(w, r, ctx)
		if err != nil {
			ctx.Logger.Info("Failed to load session",
				"err", err,
			)
			http.Error(w, "Failed to create new game", http.StatusInternalServerError)
			return
		}

		ctx.Logger.Info("Assigning new game to session",
			"sessionId", sess.Id,
		)
		ctx.Logger.Info("Created new game",
			"gameId", game.Id,
		)
		w.Header().Add("hx-push-url", fmt.Sprintf("/game/%s", game.Id))
		tmpl.ExecuteTemplate(w, "layout", game)
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

	return mux
}

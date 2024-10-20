package game

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	u "github.com/seanjh/war/utilhttp"
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

func renderFlip() http.Handler {
	tmpl := template.Must(template.ParseFiles(
		filepath.Join("templates", "game.html"),
		filepath.Join("templates", "player.html"),
		filepath.Join("templates", "battleground.html"),
		filepath.Join("templates", "warzone.html"),
	))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		game, ok := mustHaveGame(w, "1")
		if !ok {
			return
		}
		tmpl.ExecuteTemplate(w, "game", game)
	})
}

func LoadGame(id string) (*Game, error) {
	game, ok := games[id]
	if !ok {
		return nil, fmt.Errorf("Game not found: %s", id)
	}
	return game, nil
}

func mustHaveGame(w http.ResponseWriter, id string) (*Game, bool) {
	game, err := LoadGame("1")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

func createGame() http.Handler {
	tmpl := loadGameTemplates()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		game := NewGame()
		games[game.Id] = game
		log.Printf("Created new game: %s", game.Id)
		w.Header().Add("hx-push-url", fmt.Sprintf("/game/%s", game.Id))
		tmpl.ExecuteTemplate(w, "layout", game)
	})
}

func renderGame() http.Handler {
	tmpl := loadGameTemplates()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		game, ok := mustHaveGame(w, "1")
		if !ok {
			return
		}
		tmpl.ExecuteTemplate(w, "layout", game)
	})
}

func renderHome() http.Handler {
	tmpl := template.Must(template.ParseFiles(
		filepath.Join("templates", "layout.html"),
		filepath.Join("templates", "home.html"),
	))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "layout", nil)
	})
}

// Temporary global game instances
var games map[string]*Game

func SetupHandlers() {
	games = make(map[string]*Game)
	http.Handle("GET /", u.LogRequest(renderHome()))
	http.Handle("GET /lobby", u.LogRequest(renderHome()))
	http.Handle("POST /game", u.LogRequest(createGame()))
	http.Handle("GET /game", u.LogRequest(renderGame()))
	http.Handle("POST /flip", u.LogRequest(renderFlip()))
}

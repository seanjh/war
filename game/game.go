package game

import (
	"html/template"
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
	Player1 Player
	Player2 Player
	Battle  Battle
}

// NewGame returns a new Game with 2 Players with equal cuts of a new Deck.
func NewGame() *Game {
	deck := NewDeck()
	deck.Shuffle(NewRiffleShuffler())
	d1, d2 := deck.Cut()
	return &Game{
		Player1: *NewPlayer(d1, "1", "Player One"),
		Player2: *NewPlayer(d2, "2", "Player Two"),
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
		tmpl.ExecuteTemplate(w, "game", game)
	})
}

func renderGame() http.Handler {
	tmpl := template.Must(template.ParseFiles(
		filepath.Join("templates", "layout.html"),
		filepath.Join("templates", "game.html"),
		filepath.Join("templates", "player.html"),
		filepath.Join("templates", "battleground.html"),
		filepath.Join("templates", "warzone.html"),
	))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "game", game)
	})
}

func renderHome() http.Handler {
	tmpl := template.Must(template.ParseFiles(
		filepath.Join("templates", "layout.html"),
		filepath.Join("templates", "game.html"),
		filepath.Join("templates", "player.html"),
		filepath.Join("templates", "battleground.html"),
		filepath.Join("templates", "warzone.html"),
	))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "home", nil)
	})
}

// Temporary global game instance
var game *Game

func SetupHandlers() {
	game = NewGame()
	http.Handle("/", u.RequireReadOnlyMethods(u.LogRequest(renderHome())))
	http.Handle("/game", u.RequireReadOnlyMethods(u.LogRequest(renderGame())))
	http.Handle("/flip", u.RequireMethods(u.LogRequest(renderFlip()), http.MethodPost))
}

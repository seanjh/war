package game

import (
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/seanjh/war/utilhttp"
)

type Player struct {
	Deck     []Card
	InBattle Card
}

type Card struct {
	Slug string
}

type Game struct {
	Player1 Player
	Player2 Player
}

func NewGame() Game {
	return Game{
		Player1: Player{
			Deck:     []Card{},
			InBattle: Card{Slug: "2C"},
		},
		Player2: Player{
			Deck:     []Card{},
			InBattle: Card{Slug: "2H"},
		},
	}
}

func RenderGame() http.Handler {
	tmpl := template.Must(template.ParseFiles(
		filepath.Join("templates", "layout.html"),
		filepath.Join("templates", "game.html"),
		filepath.Join("templates", "player.html"),
		filepath.Join("templates", "battleground.html"),
	))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		game := NewGame()
		log.Printf("Rendering new game: %v", game)
		tmpl.ExecuteTemplate(w, "layout", game)
	})
}

func SetupHandlers() {
	http.Handle("/", utilhttp.RequireReadOnlyMethods(utilhttp.LogRequest(RenderGame())))
}

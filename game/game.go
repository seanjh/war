package game

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	u "github.com/seanjh/war/utilhttp"
)

// Temporary global game instance
var game *Game

var standardDeck = Deck{
	Card{Suit: SuitClub, Value: Two},
	Card{},
}

type Deck []Card

type Suit string
type FaceValue int

type Card struct {
	Suit  Suit
	Value FaceValue
}

const (
	SuitClub    Suit = "C"
	SuitHeart   Suit = "H"
	SuitDiamond Suit = "D"
	SuitSpade   Suit = "S"
)

func (s Suit) Name() string {
	switch s {
	case SuitClub:
		return "Clubs"
	case SuitHeart:
		return "Hearts"
	case SuitDiamond:
		return "Diamonds"
	case SuitSpade:
		return "Spades"
	default:
		return "Unrecognized"
	}
}

const (
	Two   FaceValue = 2
	Three FaceValue = 3
	Four  FaceValue = 4
	Five  FaceValue = 5
	Six   FaceValue = 6
	Seven FaceValue = 7
	Eight FaceValue = 8
	Nine  FaceValue = 9
	Ten   FaceValue = 10
	Jack  FaceValue = 11
	Queen FaceValue = 12
	King  FaceValue = 13
	Ace   FaceValue = 14
)

func (v FaceValue) Slug() string {
	switch v {
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	case Ace:
		return "A"
	default:
		return fmt.Sprint(v)
	}
}

func (v FaceValue) Name() string {
	switch v {
	case Two:
		return "Two"
	case Three:
		return "Three"
	case Four:
		return "Four"
	case Five:
		return "Five"
	case Six:
		return "Six"
	case Seven:
		return "Seven"
	case Eight:
		return "Eight"
	case Nine:
		return "Nine"
	case Ten:
		return "Ten"
	case Jack:
		return "Jack"
	case Queen:
		return "Queen"
	case King:
		return "King"
	case Ace:
		return "Ace"
	default:
		return "Unrecognized"
	}
}

func (c Card) Name() string {
	return fmt.Sprintf("%s of %s", c.Value.Name(), c.Suit.Name())
}

func (c Card) Slug() string {
	return fmt.Sprintf("%s%s", c.Value.Slug(), c.Suit)
}

type Player struct {
	Deck     Deck
	InBattle Card
	Name     string
	Id       string
}

type Game struct {
	Player1 Player
	Player2 Player
}

func newShuffledDeck() Deck {
	d := make([]Card, 52)
	return d
}

func splitDeck(d Deck) (Deck, Deck) {
	return Deck{}, Deck{}
}

func newGame() *Game {
	deck := newShuffledDeck()
	deck1, deck2 := splitDeck(deck)
	return &Game{
		Player1: Player{
			Deck:     deck1,
			Id:       "1",
			InBattle: Card{Suit: SuitClub, Value: Two},
			Name:     "One",
		},
		Player2: Player{
			Deck:     deck2,
			Id:       "1",
			InBattle: Card{Suit: SuitHeart, Value: Two},
			Name:     "Two",
		},
	}
}

func renderPage() http.Handler {
	tmpl := template.Must(template.ParseFiles(
		filepath.Join("templates", "layout.html"),
		filepath.Join("templates", "game.html"),
		filepath.Join("templates", "player.html"),
		filepath.Join("templates", "battleground.html"),
		filepath.Join("templates", "warzone.html"),
	))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "layout", game)
	})
}

func flip() func(http.ResponseWriter, *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		filepath.Join("templates", "game.html"),
		filepath.Join("templates", "player.html"),
		filepath.Join("templates", "battleground.html"),
		filepath.Join("templates", "warzone.html"),
	))
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "game", game)
	}
}

func SetupHandlers() {
	game = newGame()
	http.Handle("/", u.RequireReadOnlyMethods(u.LogRequest(renderPage())))
	http.Handle("/flip", u.RequireMethods(u.LogRequest(http.HandlerFunc(flip())), http.MethodPost))
}

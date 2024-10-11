package game

import (
	"fmt"
	"html/template"
	"log"
	"math/rand/v2"
	"net/http"
	"path/filepath"

	u "github.com/seanjh/war/utilhttp"
)

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

var SuitNames = map[Suit]string{
	SuitClub:    "Clubs",
	SuitDiamond: "Diamonds",
	SuitHeart:   "Hearts",
	SuitSpade:   "Spades",
}

func (s Suit) Name() string {
	name, ok := SuitNames[s]
	if !ok {
		return ""
	}
	return name
}

const (
	Jack  FaceValue = 11
	Queen FaceValue = 12
	King  FaceValue = 13
	Ace   FaceValue = 14
)

var FaceValueSlugs = map[FaceValue]string{
	Jack:  "J",
	Queen: "Q",
	King:  "K",
	Ace:   "A",
}

func (v FaceValue) Slug() string {
	if v > Ace {
		return ""
	}
	slug, ok := FaceValueSlugs[v]
	if !ok {
		return fmt.Sprint(v)
	}
	return slug
}

var FaceValueNames = map[FaceValue]string{
	2:     "Two",
	3:     "Three",
	4:     "Four",
	5:     "Five",
	6:     "Six",
	7:     "Seven",
	8:     "Eight",
	9:     "Nine",
	10:    "Ten",
	Jack:  "Jack",
	Queen: "Queen",
	King:  "King",
	Ace:   "Ace",
}

func (v FaceValue) Name() string {
	name, ok := FaceValueNames[v]
	if !ok {
		return "N/A"
	}
	return name
}

func (c Card) Name() string {
	return fmt.Sprintf("%s of %s", c.Value.Name(), c.Suit.Name())
}

func (c Card) Slug() string {
	return fmt.Sprintf("%s%s", c.Value.Slug(), c.Suit)
}

type Deck []Card

func NewDeck() Deck {
	d := make([]Card, 0)
	for _, s := range []Suit{SuitClub, SuitDiamond, SuitHeart, SuitSpade} {
		var v FaceValue
		for v = 2; v <= Ace; v++ {
			d = append(d, Card{Suit: s, Value: v})
		}
	}
	return d
}

// cutDeck returns 2 new decks, each containing exactly 1/2 of the original deck, with
// the extra card (in odd-sized decks) added to the first (left) deck.
func (d Deck) cut() (Deck, Deck) {
	left, right := make(Deck, 0), make(Deck, 0)
	for i := 0; i < len(d); i++ {
		c := Card{d[i].Suit, d[i].Value}
		if i&1 == 1 {
			left = append(left, c)
		} else {
			right = append(right, c)
		}
	}
	return left, right
}

type Shuffler = func(d Deck) Deck

// riffleShuffler returns a copy of the deck using a rough approximation of the "Riffle shuffle"
// technique - where cards are cut into 2 smaller decks, and interleaved. See
// [Riffle shuffle permutation] for details.
//
// [Riffle shuffle permutation]: https://en.wikipedia.org/wiki/Riffle_shuffle_permutation
func riffleShuffler(d Deck) Deck {
	result := make(Deck, 0)

	left, right := d.cut()
	log.Printf("Cut deck into 2 packages. left=%d, right=%d", len(left), len(right))

	leftI := 0
	rightI := 0
	for i := 0; i < len(left)+len(right); i++ {
		// interleave cards randomly from each cut
		takeLeft := rand.IntN(1) == 1
		if takeLeft && leftI < len(left) {
			result = append(result, left[leftI])
			leftI++
		} else if rightI < len(right) {
			result = append(result, right[rightI])
			rightI++
		} else {
			log.Printf("Cannot take from left (len: %d) or right (len: %d) cuts: %d", len(left), len(right), i)
			break
		}
	}
	return result
}

// It takes just seven ordinary, imperfect shuffles to mix a deck of cards
// thoroughly, researchers have found. Fewer are not enough and more do not
// significantly improve the mixing.
//
// [In Shuffling Cards, 7 Is Winning Number]: https://www.nytimes.com/1990/01/09/science/in-shuffling-cards-7-is-winning-number.html
const defaultShuffleRounds = 7

func (d *Deck) shuffle(shuffler Shuffler) {
	log.Printf("Performing shuffle for deck. size=%d, rounds=%d", len(*d), defaultShuffleRounds)
	for i := 0; i < defaultShuffleRounds; i++ {
		*d = shuffler(*d)
		log.Printf("Finished riffle shuffle round #%d", i+1)
	}
}

type Player struct {
	Deck     Deck
	InBattle Card
	Name     string
	Id       string
}

// NewPlayer returns a new player.
func NewPlayer(d Deck, id string, name string) *Player {
	p := Player{Deck: d, Id: id, Name: name}
	return &p
}

type Battle struct {
	Battle  map[string]Card
	Warzone map[string][]Card
}

type Game struct {
	Player1 Player
	Player2 Player
	Battle  []Card
}

func NewGame() *Game {
	deck := NewDeck()
	deck.shuffle(riffleShuffler)
	p1d, p2d := deck.cut()
	return &Game{
		Player1: Player{
			Deck: p1d,
			Id:   "1",
			Name: "One",
		},
		Player2: Player{
			Deck: p2d,
			Id:   "2",
			Name: "Two",
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

// Temporary global game instance
var game *Game

func SetupHandlers() {
	game = NewGame()
	http.Handle("/", u.RequireReadOnlyMethods(u.LogRequest(renderPage())))
	http.Handle("/flip", u.RequireMethods(u.LogRequest(http.HandlerFunc(flip())), http.MethodPost))
}

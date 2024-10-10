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

// TODO(sean) - takes a deck
// returns a card from next()
type Shuffler interface {
    done() bool
	next() (Card, bool)
}

type RiffleShuffler struct {
    leftCut []Card
    rightCut []Card
}

func NewRiffleShuffler(c []Card) RiffleShuffler {
    left, right := cutDeck(c)
    return RiffleShuffler{left, right}
}

func (s Shuffler) next() (Card, bool) {
}

type Deck struct {
	cards    []Card
	shuffler Shuffler
}

func newDeck() Deck {
	c := make([]Card, 0)
	d := Deck{c}
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
func cutDeck(d Deck) (Deck, Deck) {
	left, right := make(Deck, 0), make(Deck, 0)
	for i := 0; i < len(d); i++ {
		c := Card{d[i].Suit, d[i].Value}
		if i%2 == 0 {
			left = append(left, c)
		} else {
			right = append(right, c)
		}
	}
	return left, right
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

// riffleShuffle returns a copy of the deck using a rough approximation of the "Riffle shuffle"
// technique - where cards are cut into 2 smaller decks, and interleaved. See
// [Riffle shuffle permutation] for details.

// IT takes just seven ordinary, imperfect shuffles to mix a deck of cards
// thoroughly, researchers have found. Fewer are not enough and more do not
// significantly improve the mixing.

//
// [In Shuffling Cards, 7 Is Winning Number]: https://www.nytimes.com/1990/01/09/science/in-shuffling-cards-7-is-winning-number.html
const defaultShuffleRounds = 7

//
// [Riffle shuffle permutation]: https://en.wikipedia.org/wiki/Riffle_shuffle_permutation
func riffleShuffle(d Deck) Deck {
	log.Printf("Performing riffle shuffle for deck. size=%d, rounds=%d", len(d), defaultShuffleRounds)
    r := {cards: make(Deck, 0, len(d)), shuffler:  {}}
	for i, c := range d {
		r[i] = Card{c.Suit, c.Value}
	}

	for i := 0; i < defaultShuffleRounds; i++ {
		left, right := cutDeck(r)
		log.Printf("Cut deck into 2 packages. left=%d, right=%d", len(left.cards), len(right.cards))

		leftI := 0
		rightI := 0
		for j := 0; j < len(left.cards)+len(right.cards); j++ {
			// interleave cards randomly from each cut
			if rand.IntN(1) == 1 && rightI < len(right.cards) {
				r[j] = right.cards[rightI]
				rightI++
			} else {
				r[j] = left[leftI]
				leftI++
			}
		}
		log.Printf("Finished riffle shuffle round #%d", i+1)
	}
	return r
}

func newGame() *Game {
	deck := riffleShuffle(newDeck())
	deck1, deck2 := cutDeck(deck)
	return &Game{
		Player1: Player{
			Deck:     deck1,
			Id:       "1",
			InBattle: Card{Suit: SuitClub, Value: 2},
			Name:     "One",
		},
		Player2: Player{
			Deck:     deck2,
			Id:       "1",
			InBattle: Card{Suit: SuitHeart, Value: 2},
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

// Temporary global game instance
var game *Game

func SetupHandlers() {
	game = newGame()
	http.Handle("/", u.RequireReadOnlyMethods(u.LogRequest(renderPage())))
	http.Handle("/flip", u.RequireMethods(u.LogRequest(http.HandlerFunc(flip())), http.MethodPost))
}

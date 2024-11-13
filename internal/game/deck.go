package game

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
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

// ConvertCardSlug converts a card slug like "10C" or "AD" into the corresponding Card.
func ConvertCardSlug(s string) (Card, error) {
	suit := s[len(s)-1:]
	rawValue := s[:len(s)-1]
	var value FaceValue
	switch rawValue {
	case "A":
		value = Ace
	case "J":
		value = Jack
	case "Q":
		value = Queen
	case "K":
		value = King
	default:
		val, err := strconv.Atoi(rawValue)
		if err != nil {
			return Card{}, fmt.Errorf("invalid card slug: %s", s)
		}
		value = FaceValue(val)
	}
	return Card{Suit: Suit(suit), Value: FaceValue(value)}, nil
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

// ConvertDeck converts a comma-separated string of card slugs into a Deck.
func ConvertDeck(s string) (Deck, error) {
	slugs := strings.Split(s, ",")
	d := make([]Card, 0)
	for _, slug := range slugs {
		card, err := ConvertCardSlug(slug)
		if err != nil {
			return nil, fmt.Errorf("invalid card %s: %w", slug, err)
		}
		d = append(d, card)
	}
	return d, nil
}

// Cut returns 2 new decks, each containing exactly 1/2 of the original deck, with
// the extra card (in odd-sized decks) added to the first (left) deck.
func (d Deck) Cut() (Deck, Deck) {
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

// Return the serialized Deck
func (d Deck) String() string {
	r := make([]string, len(d))
	for i, card := range d {
		r[i] = card.Slug()
	}
	return strings.Join(r, ",")
}

type Shuffler interface {
	Shuffle(Deck) Deck
}

// RiffleShuffler is a shuffler that uses a rough approximation of the "Riffle shuffle"
// technique - where cards are cut into 2 smaller decks, and interleaved. See
// [Riffle shuffle permutation] for details.
//
// [Riffle shuffle permutation]: https://en.wikipedia.org/wiki/Riffle_shuffle_permutation
type RiffleShuffler struct {
	// random returns a value in the range [0.0,1.0), which determines from
	// which cut to pull the next card during a shuffle.
	random func() float32
}

func NewRiffleShuffler() *RiffleShuffler {
	s := RiffleShuffler{random: rand.Float32}
	return &s
}

func (s RiffleShuffler) Shuffle(d Deck) Deck {
	r := make(Deck, 0)

	left, right := d.Cut()
	li := 0
	ri := 0

	for i := 0; i < len(left)+len(right); i++ {
		leftRemain := li < len(left)
		rightRemain := ri < len(right)
		leftPreferred := s.random() < 0.5

		if leftRemain && !rightRemain {
			r = append(r, left[li])
			li++
		} else if rightRemain && !leftRemain {
			r = append(r, right[ri])
			ri++
		} else if leftPreferred {
			r = append(r, left[li])
			li++
		} else {
			r = append(r, right[ri])
			ri++
		}
	}
	return r
}

// It takes just seven ordinary, imperfect shuffles to mix a deck of cards
// thoroughly, researchers have found. Fewer are not enough and more do not
// significantly improve the mixing.
//
// [In Shuffling Cards, 7 Is Winning Number]: https://www.nytimes.com/1990/01/09/science/in-shuffling-cards-7-is-winning-number.html
const defaultShuffleRounds = 7

// Shuffle randomly mixes the cards in the deck with the given shuffler.
func (d *Deck) Shuffle(s Shuffler) {
	for i := 0; i < defaultShuffleRounds; i++ {
		*d = s.Shuffle(*d)
	}
}

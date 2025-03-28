package game

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardName(t *testing.T) {
	testCases := []struct {
		card     Card
		expected string
	}{
		{
			card:     Card{Suit: "C", Value: 2},
			expected: "Two of Clubs",
		},
		{
			card:     Card{Suit: "H", Value: 3},
			expected: "Three of Hearts",
		},
		{
			card:     Card{Suit: "D", Value: 11},
			expected: "Jack of Diamonds",
		},
		{
			card:     Card{Suit: "S", Value: 14},
			expected: "Ace of Spades",
		},
	}

	for _, c := range testCases {
		t.Run("card human-readable name", func(t *testing.T) {
			assert.Equal(t, c.card.Name(), c.expected)
		})
	}
}

func TestCardSlug(t *testing.T) {
	testCases := []struct {
		card     Card
		expected string
	}{
		{
			card:     Card{Suit: "C", Value: 2},
			expected: "2C",
		},
		{
			card:     Card{Suit: "H", Value: 3},
			expected: "3H",
		},
		{
			card:     Card{Suit: "D", Value: 11},
			expected: "JD",
		},
		{
			card:     Card{Suit: "S", Value: 14},
			expected: "AS",
		},
	}
	assert := assert.New(t)
	for _, c := range testCases {
		t.Run("convert card to slug-style string", func(t *testing.T) {
			assert.Equal(c.card.Slug(), c.expected)
		})
	}
}

func TestNewDeck(t *testing.T) {
	expected := Deck{
		Card{Suit: "C", Value: 2},
		Card{Suit: "C", Value: 3},
		Card{Suit: "C", Value: 4},
		Card{Suit: "C", Value: 5},
		Card{Suit: "C", Value: 6},
		Card{Suit: "C", Value: 7},
		Card{Suit: "C", Value: 8},
		Card{Suit: "C", Value: 9},
		Card{Suit: "C", Value: 10},
		Card{Suit: "C", Value: Jack},
		Card{Suit: "C", Value: Queen},
		Card{Suit: "C", Value: King},
		Card{Suit: "C", Value: Ace},
		Card{Suit: "D", Value: 2},
		Card{Suit: "D", Value: 3},
		Card{Suit: "D", Value: 4},
		Card{Suit: "D", Value: 5},
		Card{Suit: "D", Value: 6},
		Card{Suit: "D", Value: 7},
		Card{Suit: "D", Value: 8},
		Card{Suit: "D", Value: 9},
		Card{Suit: "D", Value: 10},
		Card{Suit: "D", Value: Jack},
		Card{Suit: "D", Value: Queen},
		Card{Suit: "D", Value: King},
		Card{Suit: "D", Value: Ace},
		Card{Suit: "H", Value: 2},
		Card{Suit: "H", Value: 3},
		Card{Suit: "H", Value: 4},
		Card{Suit: "H", Value: 5},
		Card{Suit: "H", Value: 6},
		Card{Suit: "H", Value: 7},
		Card{Suit: "H", Value: 8},
		Card{Suit: "H", Value: 9},
		Card{Suit: "H", Value: 10},
		Card{Suit: "H", Value: Jack},
		Card{Suit: "H", Value: Queen},
		Card{Suit: "H", Value: King},
		Card{Suit: "H", Value: Ace},
		Card{Suit: "S", Value: 2},
		Card{Suit: "S", Value: 3},
		Card{Suit: "S", Value: 4},
		Card{Suit: "S", Value: 5},
		Card{Suit: "S", Value: 6},
		Card{Suit: "S", Value: 7},
		Card{Suit: "S", Value: 8},
		Card{Suit: "S", Value: 9},
		Card{Suit: "S", Value: 10},
		Card{Suit: "S", Value: Jack},
		Card{Suit: "S", Value: Queen},
		Card{Suit: "S", Value: King},
		Card{Suit: "S", Value: Ace},
	}

	assert.Equal(t, expected, NewDeck())
}

func TestCutDeck(t *testing.T) {
	testCases := []struct {
		scenario      string
		deck          Deck
		expectedLeft  Deck
		expectedRight Deck
	}{
		{
			"empty deck",
			Deck{},
			Deck{},
			Deck{},
		},
		{
			"one card",
			Deck{Card{"C", 2}},
			Deck{},
			Deck{Card{"C", 2}},
		},
		{
			"two cards",
			Deck{Card{"C", 2}, Card{"H", 2}},
			Deck{Card{"H", 2}},
			Deck{Card{"C", 2}},
		},
	}

	assert := assert.New(t)
	for _, c := range testCases {
		t.Run(c.scenario, func(t *testing.T) {
			left, right := c.deck.Cut()
			assert.Equal(c.expectedLeft, left)
			assert.Equal(c.expectedRight, right)
		})
	}
}

func TestShuffle(t *testing.T) {
	testCases := []struct {
		scenario string
		deck     Deck
		expected Deck
	}{
		{
			scenario: "empty deck",
			deck:     Deck{},
			expected: Deck{},
		},
		{
			scenario: "one card",
			deck:     Deck{Card{"C", 2}},
			expected: Deck{Card{"C", 2}},
		},
		{
			scenario: "five cards",
			deck: Deck{
				Card{"C", 2},
				Card{"D", 2},
				Card{"H", 2},
				Card{"S", 2},
				Card{"S", Ace},
			},
			expected: Deck{
				Card{"H", 2},
				Card{"C", 2},
				Card{"S", 2},
				Card{"D", 2},
				Card{"S", Ace},
			},
		},
		{
			scenario: "full deck",
			deck: Deck{
				Card{Suit: "C", Value: 2},
				Card{Suit: "C", Value: 3},
				Card{Suit: "C", Value: 4},
				Card{Suit: "C", Value: 5},
				Card{Suit: "C", Value: 6},
				Card{Suit: "C", Value: 7},
				Card{Suit: "C", Value: 8},
				Card{Suit: "C", Value: 9},
				Card{Suit: "C", Value: 10},
				Card{Suit: "C", Value: Jack},
				Card{Suit: "C", Value: Queen},
				Card{Suit: "C", Value: King},
				Card{Suit: "C", Value: Ace},
				Card{Suit: "D", Value: 2},
				Card{Suit: "D", Value: 3},
				Card{Suit: "D", Value: 4},
				Card{Suit: "D", Value: 5},
				Card{Suit: "D", Value: 6},
				Card{Suit: "D", Value: 7},
				Card{Suit: "D", Value: 8},
				Card{Suit: "D", Value: 9},
				Card{Suit: "D", Value: 10},
				Card{Suit: "D", Value: Jack},
				Card{Suit: "D", Value: Queen},
				Card{Suit: "D", Value: King},
				Card{Suit: "D", Value: Ace},
				Card{Suit: "H", Value: 2},
				Card{Suit: "H", Value: 3},
				Card{Suit: "H", Value: 4},
				Card{Suit: "H", Value: 5},
				Card{Suit: "H", Value: 6},
				Card{Suit: "H", Value: 7},
				Card{Suit: "H", Value: 8},
				Card{Suit: "H", Value: 9},
				Card{Suit: "H", Value: 10},
				Card{Suit: "H", Value: Jack},
				Card{Suit: "H", Value: Queen},
				Card{Suit: "H", Value: King},
				Card{Suit: "H", Value: Ace},
				Card{Suit: "S", Value: 2},
				Card{Suit: "S", Value: 3},
				Card{Suit: "S", Value: 4},
				Card{Suit: "S", Value: 5},
				Card{Suit: "S", Value: 6},
				Card{Suit: "S", Value: 7},
				Card{Suit: "S", Value: 8},
				Card{Suit: "S", Value: 9},
				Card{Suit: "S", Value: 10},
				Card{Suit: "S", Value: Jack},
				Card{Suit: "S", Value: Queen},
				Card{Suit: "S", Value: King},
				Card{Suit: "S", Value: Ace},
			},
			expected: Deck{
				Card{Suit: "D", Value: 10},
				Card{Suit: "S", Value: 6},
				Card{Suit: "C", Value: Ace},
				Card{Suit: "H", Value: 10},
				Card{Suit: "C", Value: 5},
				Card{Suit: "D", Value: Ace},
				Card{Suit: "S", Value: 10},
				Card{Suit: "D", Value: 5},
				Card{Suit: "H", Value: Ace},
				Card{Suit: "C", Value: 9},
				Card{Suit: "H", Value: 5},
				Card{Suit: "S", Value: Ace},
				Card{Suit: "D", Value: 9},
				Card{Suit: "S", Value: 5},
				Card{Suit: "C", Value: King},
				Card{Suit: "H", Value: 9},
				Card{Suit: "C", Value: 4},
				Card{Suit: "D", Value: King},
				Card{Suit: "S", Value: 9},
				Card{Suit: "D", Value: 4},
				Card{Suit: "H", Value: King},
				Card{Suit: "C", Value: 8},
				Card{Suit: "H", Value: 4},
				Card{Suit: "S", Value: King},
				Card{Suit: "D", Value: 8},
				Card{Suit: "S", Value: 4},
				Card{Suit: "C", Value: Queen},
				Card{Suit: "H", Value: 8},
				Card{Suit: "C", Value: 3},
				Card{Suit: "D", Value: Queen},
				Card{Suit: "S", Value: 8},
				Card{Suit: "D", Value: 3},
				Card{Suit: "H", Value: Queen},
				Card{Suit: "C", Value: 7},
				Card{Suit: "H", Value: 3},
				Card{Suit: "S", Value: Queen},
				Card{Suit: "D", Value: 7},
				Card{Suit: "S", Value: 3},
				Card{Suit: "C", Value: Jack},
				Card{Suit: "H", Value: 7},
				Card{Suit: "C", Value: 2},
				Card{Suit: "D", Value: Jack},
				Card{Suit: "S", Value: 7},
				Card{Suit: "D", Value: 2},
				Card{Suit: "H", Value: Jack},
				Card{Suit: "C", Value: 6},
				Card{Suit: "H", Value: 2},
				Card{Suit: "S", Value: Jack},
				Card{Suit: "D", Value: 6},
				Card{Suit: "S", Value: 2},
				Card{Suit: "C", Value: 10},
				Card{Suit: "H", Value: 6},
			},
		},
	}

	nonRandom := func() float32 { return 0.0 }
	s := RiffleShuffler{random: nonRandom}

	assert := assert.New(t)
	for _, c := range testCases {
		t.Run(c.scenario, func(t *testing.T) {
			c.deck.Shuffle(s)
			assert.Len(c.deck, len(c.expected))
			assert.Equal(c.expected, c.deck)
		})
	}
}

func TestDeckString(t *testing.T) {
	testCases := []struct {
		deck     Deck
		expected string
	}{
		{
			deck:     []Card{},
			expected: "",
		},
		{
			deck:     []Card{{SuitHeart, 2}},
			expected: "2H",
		},
		{
			deck:     []Card{{SuitHeart, 2}, {SuitDiamond, 2}},
			expected: "2H,2D",
		},
	}

	assert := assert.New(t)
	for _, c := range testCases {
		t.Run("convert deck to string", func(t *testing.T) {
			assert.Equal(c.deck.String(), c.expected)
		})
	}
}

func TestConvertDeck(t *testing.T) {
	testCases := []struct {
		scenario string
		slug     string
		expected Deck
	}{
		{
			scenario: "empty deck",
			slug:     "",
			expected: Deck{},
		},
		{
			scenario: "two card deck",
			slug:     "10C,AD",
			expected: Deck{Card{"C", 10}, Card{"D", Ace}},
		},
		{
			scenario: "three card deck",
			slug:     "10C,AD,3H",
			expected: Deck{Card{"C", 10}, Card{"D", Ace}, Card{"H", 3}},
		},
	}

	assert := assert.New(t)
	for _, c := range testCases {
		t.Run(fmt.Sprintf("convert string to %s", c.scenario), func(t *testing.T) {
			d, err := ConvertDeck(c.slug)
			if assert.Nil(err) {
				assert.Equal(d, c.expected)
			}
		})
	}
}

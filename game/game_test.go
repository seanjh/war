package game

import (
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
		t.Run("card name", func(t *testing.T) {
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
	for _, c := range testCases {
		t.Run("card slug", func(t *testing.T) {
			assert.Equal(t, c.card.Slug(), c.expected)
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

	assert.Equal(t, expected, newDeck())
}

func TestCutDeck(t *testing.T) {
	testCases := []struct {
		scenario string
		deck     Deck
		expected struct {
			left  Deck
			right Deck
		}
	}{
		{
			"empty deck",
			Deck{},
			struct {
				left  Deck
				right Deck
			}{Deck{}, Deck{}},
		},
		{
			"one card",
			Deck{Card{"C", 2}},
			struct {
				left  Deck
				right Deck
			}{
				Deck{Card{"C", 2}},
				Deck{},
			},
		},
	}

	for _, c := range testCases {
		t.Run(c.scenario, func(t *testing.T) {
			left, right := cutDeck(c.deck)
			assert.Equal(t, c.expected.left, left)
			assert.Equal(t, c.expected.right, right)
		})
	}
}

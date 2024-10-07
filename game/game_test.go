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

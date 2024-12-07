package game

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlipCard(t *testing.T) {
	testCases := []struct {
		scenario string
		war      War
		expected War
	}{
		{
			scenario: "one card hand",
			war: War{
				Hand:       Deck{Card{Suit: "H", Value: 2}},
				Battling:   Deck{},
				Supporting: []Deck{},
			},
			expected: War{
				Hand:       Deck{},
				Battling:   Deck{Card{Suit: "H", Value: 2}},
				Supporting: []Deck{},
			},
		},
		{
			scenario: "war with one card hand",
			war: War{
				Hand:       Deck{Card{Suit: "H", Value: 3}},
				Battling:   Deck{Card{Suit: "H", Value: 2}},
				Supporting: []Deck{},
			},
			expected: War{
				Hand:       Deck{},
				Battling:   Deck{Card{Suit: "H", Value: 2}, Card{Suit: "H", Value: 3}},
				Supporting: []Deck{},
			},
		},
	}

	assert := assert.New(t)
	for _, c := range testCases {
		t.Run(fmt.Sprintf("handles flip %s", c.scenario), func(t *testing.T) {
			err := c.war.flip()
			assert.Nil(err)
			assert.Equal(c.expected, c.war)
		})
	}

}

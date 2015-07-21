package pansi

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEscapeState(t *testing.T) {
	FeedRune('')
	assert.Equal(t, globalLexer.State(), escape)
}

func TestInitalState(t *testing.T) {
	l := NewLexer()
	assert.Nil(t, l.Result(), "NewLexer's result is not nil")
	assert.Equal(t, l.State(), ground, "NewLexer's state is not ground")
}

func TestInvalidStrings(t *testing.T) {
	invalid := []string{"", "", " ", "OOOOOO", "[", "jj"}

	for _, s := range invalid {
		result := Parse(s)
		if result != nil {
			t.Errorf("AnsiEscapeCode %v is not nil", result.Type)
		}
	}
}

func TestColorCodes(t *testing.T) {
	codes := []string{"[38;5;49m", "[0m", "[38;5;49m", "[0m", "[38;5;48m", "[0m", "[38;5;48m", "[0m", "[38;5;48m"}
	for _, code := range codes {
		res := Parse(code)
		assert.Equal(t, res.Type, SetGraphicMode, "%T is not %T; %+v", res.Type, SetGraphicMode)
	}
}

func TestColorCode(t *testing.T) {
	type L struct {
		b byte
		s State
	}
	bytes := []L{
		L{0x1B, escape},
		L{0x5b, csiEntry},
		L{0x33, csiParam},
		L{0x38, csiParam},
		L{0x3B, csiParam},
		L{0x35, csiParam},
		L{0x3B, csiParam},
		L{0x34, csiParam},
		L{0x39, csiParam},
		L{0x6D, ground},
	}
	lexer := NewLexer()
	for _, l := range bytes {
		lexer.Feed(l.b)
		assert.Equal(t, lexer.State(), l.s, "%T is not %T", lexer.State(), l.s)
	}
	res := lexer.Result()
	assert.NotNil(t, res, "p.Result() should not be nil")
	assert.Equal(t, res.Type, SetGraphicMode, "res.Type should be SetGraphicMode")
	assert.Equal(t, res.Values, []int{38, 5, 49}, "res.Values got borked")
}

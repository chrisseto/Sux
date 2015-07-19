package pansi

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestInvalidStrings(t *testing.T) {
	invalid := []string{"", "", " ", "OOOOOO", "[", "jj"}

	for _, s := range invalid {
		result := Parse(s)
		if result != nil {
			t.Errorf("AnsiEscapeCode %v is not nil", result.Type)
		}
	}
}

func TestInitalState(t *testing.T) {
	p := NewParser()
	if p.Result() != nil {
		t.Errorf("NewParser's result is not nil")
	}
}

func TestEscapeState(t *testing.T) {
	p := NewParser()
	p.FeedRune('')
	_, ok := p.state.(*escape)
	if !ok {
		t.Errorf("%T is not %T", p.state, escape{})
	}
}

func TestColorCodes(t *testing.T) {
	codes := []string{"[38;5;49m", "[0m", "[38;5;49m", "[0m", "[38;5;48m", "[0m", "[38;5;48m", "[0m", "[38;5;48m"}
	for _, code := range codes {
		res := Parse(code)
		if res.Type != SetGraphicMode {
			t.Errorf("%T is not %T; %+v", res.Type, SetGraphicMode)
		}
	}
}

func TestColorCode(t *testing.T) {
	type L struct {
		b byte
		s State
	}
	bytes := []L{
		L{0x1B, &escape{}},
		L{0x5b, &csiEntry{}},
		L{0x33, &csiParam{}},
		L{0x38, &csiParam{}},
		L{0x3B, &csiParam{}},
		L{0x35, &csiParam{}},
		L{0x3B, &csiParam{}},
		L{0x34, &csiParam{}},
		L{0x39, &csiParam{}},
		L{0x6D, nil},
	}
	p := NewParser()
	for _, l := range bytes {
		p.Feed(l.b)
		if reflect.TypeOf(p.state) != reflect.TypeOf(l.s) {
			t.Errorf("%T is not %T", p.state, escape{})
		}
	}
	res := p.Result()
	assert.NotNil(t, res, "p.Result() should not be nil")
	assert.Equal(t, res.Type, SetGraphicMode, "res.Type should be SetGraphicMode")
	assert.Equal(t, res.Values, []int{38, 5, 49}, "res.Values got borked")
}

func TestColorCodeString(t *testing.T) {
	res := Parse("[38;5;49m")
	assert.NotNil(t, res, "p.Result() should not be nil")
	assert.Equal(t, res.Type, SetGraphicMode, "res.Type should be SetGraphicMode")
	assert.Equal(t, res.Values, []int{38, 5, 49}, "res.Values got borked")
}

func BenchmarkEscapeState(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := NewParser()
		p.Feed('')
	}
}

func BenchmarkStringColorParsing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Parse("[38;5;49m")
	}
}

func BenchmarkRuneColorParsing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := NewParser()
		p.FeedRune('')
		p.FeedRune('[')
		p.FeedRune('3')
		p.FeedRune('8')
		p.FeedRune(';')
		p.FeedRune('5')
		p.FeedRune(';')
		p.FeedRune('4')
		p.FeedRune('9')
		p.FeedRune('m')
		p.Result()
	}
}

func BenchmarkByteColorParsing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := NewParser()
		p.Feed(0x1B)
		p.Feed(0x5B)
		p.Feed(0x33)
		p.Feed(0x38)
		p.Feed(0x3B)
		p.Feed(0x35)
		p.Feed(0x3B)
		p.Feed(0x34)
		p.Feed(0x39)
		p.Feed(0x6D)
		p.Result()
	}
}

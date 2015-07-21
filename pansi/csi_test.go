package pansi

// import (
// 	"github.com/stretchr/testify/assert"
// 	"reflect"
// 	"testing"
// )

// func TestColorCodeString(t *testing.T) {
// 	res := Parse("[38;5;49m")
// 	assert.NotNil(t, res, "p.Result() should not be nil")
// 	assert.Equal(t, res.Type, SetGraphicMode, "res.Type should be SetGraphicMode")
// 	assert.Equal(t, res.Values, []int{38, 5, 49}, "res.Values got borked")
// }

// func TestEraseLine(t *testing.T) {
// 	res := Parse("[K")
// 	assert.NotNil(t, res, "p.Result() should not be nil")
// 	assert.Equal(t, res.Type, EraseLine, "res.Type should be EraseLine")
// }

// func TestTypes(t *testing.T) {
// 	type L struct {
// 		s string
// 		t AnsiEscapeType
// 	}
// 	inputs := []L{
// 		L{"[2J", EraseDisplay},
// 	}

// 	for _, input := range inputs {
// 		res := Parse(input.s)
// 		assert.NotNil(t, res, "%#v caused a nil Result", input)
// 		assert.Equal(t, res.Type, input.t, "Result.Type should be %v", input.t)
// 	}

// }

// func BenchmarkEscapeState(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		p := NewParser()
// 		p.Feed('')
// 	}
// }

// func BenchmarkStringColorParsing(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		Parse("[38;5;49m")
// 	}
// }

// func BenchmarkRuneColorParsing(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		p := NewParser()
// 		p.FeedRune('')
// 		p.FeedRune('[')
// 		p.FeedRune('3')
// 		p.FeedRune('8')
// 		p.FeedRune(';')
// 		p.FeedRune('5')
// 		p.FeedRune(';')
// 		p.FeedRune('4')
// 		p.FeedRune('9')
// 		p.FeedRune('m')
// 		p.Result()
// 	}
// }

// func BenchmarkByteColorParsing(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		p := NewParser()
// 		p.Feed(0x1B)
// 		p.Feed(0x5B)
// 		p.Feed(0x33)
// 		p.Feed(0x38)
// 		p.Feed(0x3B)
// 		p.Feed(0x35)
// 		p.Feed(0x3B)
// 		p.Feed(0x34)
// 		p.Feed(0x39)
// 		p.Feed(0x6D)
// 		p.Result()
// 	}
// }

package pane

import (
	"github.com/nsf/termbox-go"
	"testing"
)

func TestAppend(t *testing.T) {
	ring := NewRingBuffer(1, 10)

	for i := 0; i < 100; i++ {
		ring.Append([]termbox.Cell{
			termbox.Cell{rune(i), 0x0, 0x0},
		})
		if ring.Get(9)[0].Ch != rune(i) {
			t.Error("Expected ring.Get(9) to be", rune(i), ". Got", ring.Get(9)[0].Ch)
		}
	}
}

func TestAppendShiftsHead(t *testing.T) {
	ring := NewRingBuffer(1, 10)

	for i := 0; i < 10; i++ {
		ring.Get(i)[0].Ch = rune(i)
	}

	for i := 10; i < 100; i++ {
		if ring.Get(0)[0].Ch != rune(i-10) {
			t.Error("Expected ring.Get(0) to be", rune(i-10), ". Got", ring.Get(0)[0].Ch)
		}
		ring.Append([]termbox.Cell{
			termbox.Cell{rune(i), 0x0, 0x0},
		})
	}
}

func TestRange(t *testing.T) {
	ring := NewRingBuffer(1, 10)

	for i := 0; i < 10; i++ {
		ring.Get(i)[0].Ch = rune(i)
	}

	for i := 0; i < 100; i++ {
		tmp := ring.Range(0, 10)
		for j := 0; j < 10; j++ {
			if tmp[j][0].Ch != rune(i+j) {
				t.Error("Expected ring.Range(0, 10)[", j, "] to be", rune(i+j), ". Got", tmp[j][0].Ch)
			}
		}

		ring.Append([]termbox.Cell{
			termbox.Cell{rune(i + 10), 0x0, 0x0},
		})
	}
}

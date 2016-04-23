package pane

import (
	"github.com/nsf/termbox-go"
	"testing"
)

func TestAppend(t *testing.T) {
	ring := NewRingBuffer(make([][]termbox.Cell, 0, 10))

	for i := 0; i < 10; i++ {
		if ring.Length() != i {
			t.Error("Expected ring.Length() to be", i, ". Got", ring.Length())
		}
		ring.Append(make([]termbox.Cell, 1))
	}

	for i := 0; i < 100; i++ {
		ring.Append([]termbox.Cell{
			termbox.Cell{rune(i), 0x0, 0x0},
		})
		t.Logf("Offset for 9 is %d", ring.offset(9))

		if ring.Get(9)[0].Ch != rune(i) {
			t.Error("Expected ring.Get(9) to be", rune(i), ". Got", ring.Get(9)[0].Ch)
		}
	}

	if ring.Length() != 10 {
		t.Error("Expected ring.Length() to be 10 . Got", ring.Length())
	}
}

func TestAppendShiftsHead(t *testing.T) {
	ring := NewRingBuffer(make([][]termbox.Cell, 0, 10))

	for i := 0; i < 10; i++ {
		ring.Append([]termbox.Cell{
			termbox.Cell{rune(i), 0x0, 0x0},
		})
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
	ring := NewRingBuffer(make([][]termbox.Cell, 0, 10))

	for i := 0; i < 10; i++ {
		ring.Append([]termbox.Cell{
			termbox.Cell{rune(i), 0x0, 0x0},
		})
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

func TestClear(t *testing.T) {
	ring := NewRingBuffer(make([][]termbox.Cell, 0, 10))

	for i := 0; i < 10; i++ {
		ring.Append([]termbox.Cell{
			termbox.Cell{rune(i), 0x0, 0x0},
		})
	}

	for i := 0; i < 100; i++ {
		ring.Clear()
		if ring.Length() != 0 {
			t.Error("Expected ring.Length() to be 1. Got", ring.Length())
		}
		ring.Append([]termbox.Cell{
			termbox.Cell{rune(i), 0x0, 0x0},
		})

		for j := 1; j < 5; j++ {
			if ring.Length() != j {
				t.Errorf("Expected ring.Length() to be %d. Got %d", j, ring.Length())
			}
			ring.Append([]termbox.Cell{
				termbox.Cell{rune(0), 0x0, 0x0},
			})
		}

		if ring.Get(0)[0].Ch != rune(i) {
			t.Error("Expected ring.Get(0) to be", rune(i), ". Got", ring.Get(0)[0].Ch)
		}
	}
}

func TestClearPreserves(t *testing.T) {
	ring := NewRingBuffer(make([][]termbox.Cell, 0, 10))

	for i := 0; i < 10; i++ {
		ring.Append([]termbox.Cell{
			termbox.Cell{rune(i), 0x0, 0x0},
		})
	}

	ring.Clear()

	if ring.Length() != 0 {
		t.Errorf("Expected ring.Length() to be 0. Got %d", ring.Length())
	}

	for i := 0; i < 10; i++ {
		if ring.buffer[i][0].Ch != rune(i) {
			t.Errorf("Expected ring.buffer[%d][0].Ch to be %s. Got %s", i, string(i), string(ring.buffer[i][0].Ch))
		}
	}
}

func TestClearWorks(t *testing.T) {
	ring := NewRingBuffer(make([][]termbox.Cell, 0, 10))

	for i := 0; i < 100; i++ {
		ring.Clear()
		for j := 0; j < i; j++ {
			ring.Append([]termbox.Cell{
				termbox.Cell{rune(j), 0x0, 0x0},
			})

			if ring.Length() != j+1 && !(j > 9 && ring.Length() == 10) {
				t.Errorf("Expected ring.Length() to be %d or 10 got %d", j+1, ring.Length())
			}

			if ring.Get(ring.Length() - 1)[0].Ch != rune(j) {
				t.Fatalf("Expected ring.Get(%d) to be %v. Got %v", ring.Length()-1, rune(j), ring.Get(ring.Length() - 1)[0].Ch)
			}
		}
	}
}

func TestClearWithRange(t *testing.T) {
	ring := NewRingBuffer(make([][]termbox.Cell, 0, 10))

	for i := 15; i < 16; i++ {
		ring.Clear()
		for j := 1; j < i; j++ {
			ring.Append([]termbox.Cell{
				termbox.Cell{rune(j), 0x0, 0x0},
			})

			data := ring.Range(0, j)

			if len(data) != j && !(j > 9 && len(data) == 10) {
				t.Errorf("Expected len(data) to be %d or 10. Got %d", j, len(data))
			}

			for index, cells := range data {
				expected := rune(j - (len(data) - (index + 1)))

				if cells[0].Ch != expected {
					t.Logf("i = %d j = %d index = %d len(data) = %d", i, j, index, len(data))
					t.Errorf("Expected data[%d] to be %v. Got %v", index, expected, data[index][0].Ch)
				}
			}
		}
	}
}

func BenchmarkAppending(b *testing.B) {
	ring := NewRingBuffer(make([][]termbox.Cell, 0, 10))

	for i := 0; i < b.N; i++ {
		ring.Append([]termbox.Cell{
			termbox.Cell{rune(i), 0x0, 0x0},
		})
	}
}

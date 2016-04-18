package pane

import (
	"github.com/nsf/termbox-go"
	// "log"
)

type RingBuffer struct {
	buffer [][]termbox.Cell
	index  int
}

func NewRingBuffer(width, length int) RingBuffer {
	r := RingBuffer{
		buffer: make([][]termbox.Cell, length),
		index:  0,
	}

	for i := 0; i < length; i++ {
		r.buffer[i] = make([]termbox.Cell, width)
	}

	return r
}

func (r *RingBuffer) Append(data []termbox.Cell) {
	r.buffer[r.index] = data

	if r.index < len(r.buffer) {
		r.index++
	} else {
		r.index = 0
	}
}

func (r *RingBuffer) Get(i int) []termbox.Cell {
	return r.buffer[r.index+i]
}

func (r *RingBuffer) Range(start, length int) [][]termbox.Cell {
	// log.Printf("Getting Range %d-%d\n%+v", start, start+length, r.buffer[r.index+start])
	return r.buffer[r.index+start : r.index+start+length]
}

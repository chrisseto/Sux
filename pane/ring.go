package pane

import (
	"github.com/nsf/termbox-go"
	"log"
)

type RingBuffer struct {
	buffer [][]termbox.Cell
	index  int
}

func NewRingBuffer(initial [][]termbox.Cell) RingBuffer {
	r := RingBuffer{
		index:  0,
		buffer: initial,
	}

	return r
}

func (r *RingBuffer) Append(data []termbox.Cell) {
	log.Printf("Appending at index %d with length %d", r.index, len(r.buffer))
	if cap(r.buffer) != len(r.buffer) {
		r.buffer = append(r.buffer, data)
		return
	}
	r.buffer[r.index] = data

	r.index++
	if r.index >= len(r.buffer) {
		r.index = 0
	}
}

func (r *RingBuffer) Get(i int) []termbox.Cell {
	log.Printf("Getting index %d, internal index %d, resolved to %d", i, r.index, r.offset(i))
	return r.buffer[r.offset(i)]
}

func (r *RingBuffer) Set(i int, data []termbox.Cell) {
	r.buffer[r.offset(i)] = data
}

func (r *RingBuffer) offset(i int) int {
	if cap(r.buffer) != len(r.buffer) {
		return i
	}
	if r.index+i >= len(r.buffer) {
		return (r.index + i) - len(r.buffer)
	}
	return r.index + i
}

func (r *RingBuffer) Range(begin, length int) [][]termbox.Cell {
	start := r.offset(begin)
	if length > len(r.buffer) {
		length = len(r.buffer)
	}
	//Golang slices are EXCLUSIVE
	//IE [1...10][0:9] -> [1,2,3,4,5,6,7,8]
	end := r.offset(begin + length - 1)

	if end < start {
		log.Printf("Range request of %d-%d resolved to %d-%d + %d-%d", begin, begin+length, start, len(r.buffer)-1, 0, end)
		return append(r.buffer[start:len(r.buffer)], r.buffer[0:end+1]...)
	}

	log.Printf("Range request of %d-%d resolved to %d-%d", begin, begin+length, start, end)
	return r.buffer[start : end+1]
}

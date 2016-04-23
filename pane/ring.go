package pane

import (
	"github.com/nsf/termbox-go"
	"log"
)

type RingBuffer struct {
	len    int
	index  int
	buffer [][]termbox.Cell
}

func NewRingBuffer(initial [][]termbox.Cell) RingBuffer {
	r := RingBuffer{
		len:    len(initial),
		index:  0,
		buffer: initial,
	}

	return r
}

func (r *RingBuffer) Length() int {
	return r.len
}

func (r *RingBuffer) Clear() {
	r.index++
	r.len = 0
}

func (r *RingBuffer) Append(data []termbox.Cell) {
	if cap(r.buffer) != len(r.buffer) {
		r.buffer = append(r.buffer, data)
	} else {
		r.buffer[r.offset(r.len)] = data
	}

	if r.len < cap(r.buffer) {
		r.len++
	} else {
		r.index++
		if r.index >= cap(r.buffer) {
			r.index = 0
		}
	}
}

func (r *RingBuffer) Get(i int) []termbox.Cell {
	return r.buffer[r.offset(i)]
}

func (r *RingBuffer) Set(i int, data []termbox.Cell) {
	r.buffer[r.offset(i)] = data
}

func (r *RingBuffer) offset(i int) int {
	return (r.index + i) % cap(r.buffer)
}

func (r *RingBuffer) Range(begin, length int) [][]termbox.Cell {
	start := r.offset(begin)
	if length > r.Length() {
		length = r.Length()
	}

	//Golang slices are EXCLUSIVE
	//IE [1...10][0:9] -> [1,2,3,4,5,6,7,8]
	end := r.offset(begin + length - 1)
	log.Printf("Range(%d, %d) -> [%d : %d]", begin, length, start, end)

	if end < start {
		return append(r.buffer[start:len(r.buffer)], r.buffer[0:end+1]...)
	}

	return r.buffer[start : end+1]
}

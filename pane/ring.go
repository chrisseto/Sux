package pane

import (
	"github.com/nsf/termbox-go"
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
	r.len = 0
	if cap(r.buffer) == len(r.buffer) {
		r.index++
		if r.index >= cap(r.buffer) {
			r.index = 0
		}
	} else {
		r.index = len(r.buffer)
	}
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
	if length == 0 || r.Length() == 0 {
		return r.buffer[0:0]
	}

	start := r.offset(begin)
	if length > r.Length() {
		length = r.Length()
	}

	//Golang slices are EXCLUSIVE
	//IE [1...10][0:9] -> [1,2,3,4,5,6,7,8]
	end := r.offset(begin + length - 1)

	if end < start {
		return append(r.buffer[start:len(r.buffer)], r.buffer[0:end+1]...)
	}

	return r.buffer[start : end+1]
}

//TODO Test me
func (r *RingBuffer) RollBack(count int) {
	r.index -= count
	if r.index < 0 {
		r.index = len(r.buffer) + r.index
	}
}

func (r *RingBuffer) Tail(count int) [][]termbox.Cell {
	if count > r.Length() {
		count = r.Length()
	}

	return r.Range(r.Length()-count, count)
}

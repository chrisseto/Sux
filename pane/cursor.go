package pane

type Cursor struct {
	x, y          int
	width, height int
}

func clamp(val, lo, hi int) int {
	if val < lo {
		return lo
	}

	if val > hi {
		return hi
	}

	return val
}

func NewCursor(width, height int) *Cursor {
	return &Cursor{x: 0, y: 0, width: width, height: height}
}

func (c *Cursor) Up(count int) {
	c.y = clamp(c.y-count, 0, c.height)
}

func (c *Cursor) Down(count int) {
	c.Up(-count)
}

func (c *Cursor) Right(count int) {
	c.x = clamp(c.x+count, 0, c.width)
}

func (c *Cursor) Left(count int) {
	c.Right(-count)
}

func (c *Cursor) SetX(x int) {
	c.x = clamp(x, 0, c.width)
}

func (c *Cursor) SetY(y int) {
	c.y = clamp(y, 0, c.height)
}

func (c *Cursor) Set(x, y int) {
	c.x = clamp(x, 0, c.width)
	c.y = clamp(y, 0, c.height)
}

func (c *Cursor) Get() (int, int) {
	return c.x, c.y
}

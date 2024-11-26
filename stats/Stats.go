package stats

import "fmt"

type Stats struct {
	User string
	Down int
	Up   int
}

func (s Stats) ToString() string {
	return fmt.Sprintf("↑ %d (mb)\n↓ %d (mb)", s.Up, s.Down)
}

func (s Stats) ToOneLineString() string {
	return fmt.Sprintf("↓ %d (mb) ↑ %d (mb)", s.Down, s.Up)
}

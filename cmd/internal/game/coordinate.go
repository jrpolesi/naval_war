package game

import "math/rand/v2"

type Coordinate struct {
	X int
	Y int
}

func newRandomCoordinate(maxX, maxY int) Coordinate {
	return Coordinate{
		X: rand.IntN(maxX),
		Y: rand.IntN(maxY),
	}
}

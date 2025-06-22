package game

type GameMap struct {
	Ships []*Ship
	Size  [2]int
}

func newGameMap(sizeX, sizeY int) GameMap {
	return GameMap{
		Ships: make([]*Ship, 0),
		Size:  [2]int{sizeX, sizeY},
	}
}

func (gm *GameMap) AddShip(ship *Ship) error {
	gm.Ships = append(gm.Ships, ship)
	return nil
}

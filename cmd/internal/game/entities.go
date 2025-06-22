package game

type Coordinate struct {
	X int
	Y int
}

type Team struct {
	ID      string
	Players []*Player
}

type Ship struct {
	ID string
	Position  Coordinate
	DamagedBy *Player   
}

type Map struct {
	Ships []*Ship
	Size  [2]int 
}

type CurrentRound struct {
	Actions []Action
}

type Action struct {
	Type    string       
	Payload ActionPayload
}

type ActionPayload struct {
	Position Coordinate
}

type Player struct {
	ID string
	Name string
	IsReady bool
	Ships []Ship
}

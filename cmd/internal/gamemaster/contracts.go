package gamemaster

type Player struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IsReady bool   `json:"is_ready"`
}

type Team struct {
	ID      string   `json:"id"`
	Players []Player `json:"players"`
}

type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Ship struct {
	ID        string     `json:"id"`
	Position  Coordinate `json:"position"`
	IsDamaged bool       `json:"is_damaged"`
	DamagedBy *Player    `json:"damaged_by,omitempty"`
	IsOwner   bool       `json:"is_owner"`
}

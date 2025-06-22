package gamemaster

type GameMap struct {
	/*
	 [width, height]
	*/
	Size [2]int `json:"size"`
}

type GameResult struct {
	Winner Team `json:"winner"`
}

type GameResponse struct {
	Map          *GameMap    `json:"map,omitempty"`
	Teams        *[]Team     `json:"teams,omitempty"`
	Ships        *[]Ship     `json:"ships,omitempty"`
	DamagedShips *[]Ship     `json:"damaged_ships,omitempty"`
	GameResult   *GameResult `json:"game_result,omitempty"`
}

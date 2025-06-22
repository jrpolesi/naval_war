package gamemaster

type UpdatePlayerIntent Player

type AttackPlayerActionIntent struct {
	Position Coordinate `json:"position"`
}

type PerformPlayerActionIntent struct {
	ActionType string                   `json:"type"`
	Payload    AttackPlayerActionIntent `json:"payload"`
}

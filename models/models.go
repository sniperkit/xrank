package models

// Team represents the entity the ranking engine ranks on.
// It could be a single player or a team with multiple players.
type Team struct {
	Id string
	Score float64
}

type Tournament struct {
	Id string
	Description string
}

// GameResult stores the result of the game.
type GameResult struct {
	Sequence int
	WinnerId string
	LoserId string
	TournamentId string
}

package models

// Simulator defines simulator servisinin dışarıya sunduğu davranışları belirtir.
type Simulator interface {
	SimulateWeek(week int) error
	SimulateAllWeeks() error
	GetCurrentStandings() ([]Team, error)
	GetAllMatches() ([]Match, error)
	GetMatchesByWeek(week int) ([]Match, error)
}

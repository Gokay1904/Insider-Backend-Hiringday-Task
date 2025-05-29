package services

import (
	"database/sql"
	"errors"
	"insider-case/models"
	"math/rand"
	"time"
)

type MatchService struct {
	DB *sql.DB
}

// NewMatchService constructor
func NewMatchService(db *sql.DB) *MatchService {
	return &MatchService{DB: db}
}

func (m *MatchService) getTeams() ([]models.Team, error) {
	rows, err := m.DB.Query("SELECT id, name, strength FROM teams")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []models.Team
	for rows.Next() {
		var t models.Team
		err := rows.Scan(&t.ID, &t.Name, &t.Strength)
		if err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}

	return teams, nil
}

func (m *MatchService) GenerateRandomMatchesForWeek(week int) error {
	teams, err := m.getTeams()
	if err != nil {
		return err
	}

	// Takım sayısı tek ise, bir takım haftayı bay geçebilir
	if len(teams)%2 != 0 {
		return errors.New("takım sayısı çift olmalı veya bay takımı eklenmeli")
	}

	// Takımları karıştır
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(teams), func(i, j int) { teams[i], teams[j] = teams[j], teams[i] })

	// Maçları oluştur: 0-1, 2-3, 4-5 şeklinde eşleştir
	for i := 0; i < len(teams); i += 2 {
		homeTeam := teams[i]
		awayTeam := teams[i+1]

		// Güçlere göre skorları simüle et
		homeGoals, awayGoals := m.SimulateMatch(homeTeam.Strength, awayTeam.Strength)

		// Maçı DB'ye ekle
		err := m.CreateMatch(homeTeam.ID, awayTeam.ID, week, homeGoals, awayGoals)
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateMatch inserts a new match record into DB
func (m *MatchService) CreateMatch(homeTeamID, awayTeamID, week, homeGoals, awayGoals int) error {
	result := "Draw"
	switch {
	case homeGoals > awayGoals:
		result = "HomeWin"
	case homeGoals < awayGoals:
		result = "AwayWin"
	}

	_, err := m.DB.Exec(`
		INSERT INTO matches (home_team_id, away_team_id, week, home_goals, away_goals, result)
		VALUES (?, ?, ?, ?, ?, ?)`,
		homeTeamID, awayTeamID, week, homeGoals, awayGoals, result)
	return err
}

// GetMatchesByWeek returns matches of a given week
func (m *MatchService) GetMatchesByWeek(week int) ([]models.Match, error) {
	rows, err := m.DB.Query(`
		SELECT id, home_team_id, away_team_id, week, home_goals, away_goals
		FROM matches WHERE week = ?`, week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []models.Match
	for rows.Next() {
		var match models.Match
		err := rows.Scan(&match.ID, &match.HomeTeamID, &match.AwayTeamID, &match.Week, &match.HomeGoals, &match.AwayGoals)
		if err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}

	return matches, nil
}

func (s *SimulatorService) PredictMatchOutcome(homeID, awayID int) (float64, float64, float64, error) {
	var home models.Team
	var away models.Team

	err := s.DB.QueryRow("SELECT played, won, drawn FROM teams WHERE id = ?", homeID).
		Scan(&home.Played, &home.Won, &home.Drawn)
	if err != nil {
		return 0, 0, 0, err
	}

	err = s.DB.QueryRow("SELECT played, won, drawn FROM teams WHERE id = ?", awayID).
		Scan(&away.Played, &away.Won, &away.Drawn)
	if err != nil {
		return 0, 0, 0, err
	}

	homeWinRate := 0.0
	homeDrawRate := 0.0
	if home.Played > 0 {
		homeWinRate = float64(home.Won) / float64(home.Played)
		homeDrawRate = float64(home.Drawn) / float64(home.Played)
	}

	awayWinRate := 0.0
	awayDrawRate := 0.0
	if away.Played > 0 {
		awayWinRate = float64(away.Won) / float64(away.Played)
		awayDrawRate = float64(away.Drawn) / float64(away.Played)
	}

	homeScore := homeWinRate + (awayWinRate * 0.5)
	awayScore := awayWinRate + (homeWinRate * 0.5)
	drawScore := (homeDrawRate + awayDrawRate) / 2

	total := homeScore + awayScore + drawScore
	if total == 0 {
		return 0.33, 0.34, 0.33, nil
	}

	return homeScore / total, drawScore / total, awayScore / total, nil
}

func (s *SimulatorService) GetChampionshipProbabilities() (map[int]float64, error) {
	rows, err := s.DB.Query("SELECT id, points FROM teams")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	totalPoints := 0
	pointsMap := make(map[int]int)
	for rows.Next() {
		var id, points int
		if err := rows.Scan(&id, &points); err != nil {
			return nil, err
		}
		pointsMap[id] = points
		totalPoints += points
	}

	probs := make(map[int]float64)
	teamCount := len(pointsMap)
	for id, pts := range pointsMap {
		if totalPoints > 0 {
			probs[id] = float64(pts) / float64(totalPoints)
		} else {
			probs[id] = 1.0 / float64(teamCount)
		}
	}

	return probs, nil
}

// SimulateMatch simulates a match result based on team strengths and returns scores
func (m *MatchService) SimulateMatch(homeStrength, awayStrength int) (int, int) {
	rand.Seed(time.Now().UnixNano())

	maxHomeGoals := homeStrength/10 + 2
	if maxHomeGoals > 5 {
		maxHomeGoals = 5
	}
	maxAwayGoals := awayStrength/10 + 2
	if maxAwayGoals > 5 {
		maxAwayGoals = 5
	}

	homeGoals := rand.Intn(maxHomeGoals)
	awayGoals := rand.Intn(maxAwayGoals)
	return homeGoals, awayGoals
}

// DeleteMatchesByWeek deletes matches for a given week - useful if simülasyon tekrar yapılacaksa
func (m *MatchService) DeleteMatchesByWeek(week int) error {
	_, err := m.DB.Exec("DELETE FROM matches WHERE week = ?", week)
	return err
}

package services

import (
	"database/sql"
	"fmt"
	"insider-case/models"
	"math/rand"
	"time"
)

type SimulatorService struct {
	DB *sql.DB
}

func NewSimulatorService(db *sql.DB) *SimulatorService {
	return &SimulatorService{DB: db}
}

func (s *SimulatorService) SimulateWeek(week int) error {
	teams, err := s.getTeams()
	if err != nil {
		return err
	}
	if len(teams) < 2 {
		return fmt.Errorf("not enough teams to simulate matches")
	}

	// Şampiyonluk oranlarını al
	champProbs, err := s.GetChampionshipProbabilities()
	if err != nil {
		return err
	}

	fmt.Printf("%d week predictions of championship\n", week)
	for _, team := range teams {
		prob, ok := champProbs[team.ID]
		if !ok {
			prob = 0.0
		}
		fmt.Printf("%s: %.2f\n", team.Name, prob)
	}

	fmt.Printf("\n%d week results\n", week)

	// Maçları simüle et
	rand.Shuffle(len(teams), func(i, j int) { teams[i], teams[j] = teams[j], teams[i] })

	stats := make(map[int]*TeamStats)
	for _, team := range teams {
		stats[team.ID] = &TeamStats{Team: team}
	}

	for i := 0; i < len(teams)-1; i += 2 {
		home := teams[i]
		away := teams[i+1]

		homeGoals, awayGoals := s.simulateScore(home.Strength, away.Strength)

		_, err = s.DB.Exec(`
			INSERT INTO matches (home_team_id, away_team_id, week, home_goals, away_goals)
			VALUES (?, ?, ?, ?, ?)`,
			home.ID, away.ID, week, homeGoals, awayGoals)
		if err != nil {
			return err
		}

		s.updateStats(stats[home.ID], stats[away.ID], homeGoals, awayGoals)

		fmt.Printf("%s %d - %d %s\n", home.Name, homeGoals, awayGoals, away.Name)
	}

	return nil
}

func (s *SimulatorService) simulateScore(homeStrength, awayStrength int) (int, int) {
	rand.Seed(time.Now().UnixNano())
	homeGoals := rand.Intn(homeStrength/10 + 2)
	awayGoals := rand.Intn(awayStrength/10 + 2)
	return homeGoals, awayGoals
}

func (s *SimulatorService) updateStats(home, away *TeamStats, homeGoals, awayGoals int) {
	home.Played++
	away.Played++

	home.GF += homeGoals
	home.GA += awayGoals
	home.GoalDiff = home.GF - home.GA

	away.GF += awayGoals
	away.GA += homeGoals
	away.GoalDiff = away.GF - away.GA

	switch {
	case homeGoals > awayGoals:
		home.Won++
		home.Points += 3
		away.Lost++
	case homeGoals < awayGoals:
		away.Won++
		away.Points += 3
		home.Lost++
	default:
		home.Drawn++
		away.Drawn++
		home.Points++
		away.Points++
	}
}

func (s *SimulatorService) getTeams() ([]models.Team, error) {
	rows, err := s.DB.Query("SELECT id, name, strength FROM teams")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []models.Team
	for rows.Next() {
		var t models.Team
		if err := rows.Scan(&t.ID, &t.Name, &t.Strength); err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}
	return teams, nil
}

func (s *SimulatorService) SimulateAllWeeks() error {
	// Örneğin 5 hafta simüle edelim, bu sayıyı ihtiyaçlarına göre değiştir
	const totalWeeks = 5

	for week := 1; week <= totalWeeks; week++ {
		if err := s.SimulateWeek(week); err != nil {
			return fmt.Errorf("failed to simulate week %d: %w", week, err)
		}
	}
	return nil
}

// GetAllMatches tüm maçları döner
func (s *SimulatorService) GetAllMatches() ([]models.Match, error) {
	rows, err := s.DB.Query(`
		SELECT id, week, home_team_id, away_team_id, home_goals, away_goals
		FROM matches
		ORDER BY week, id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []models.Match
	for rows.Next() {
		var m models.Match
		err := rows.Scan(&m.ID, &m.Week, &m.HomeTeamID, &m.AwayTeamID, &m.HomeGoals, &m.AwayGoals)
		if err != nil {
			return nil, err
		}
		matches = append(matches, m)
	}

	return matches, nil
}

// GetMatchesByWeek belirli haftaya ait maçları döner
func (s *SimulatorService) GetMatchesByWeek(week int) ([]models.Match, error) {
	rows, err := s.DB.Query(`
		SELECT id, week, home_team_id, away_team_id, home_goals, away_goals
		FROM matches
		WHERE week = ?
		ORDER BY id
	`, week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []models.Match
	for rows.Next() {
		var m models.Match
		err := rows.Scan(&m.ID, &m.Week, &m.HomeTeamID, &m.AwayTeamID, &m.HomeGoals, &m.AwayGoals)
		if err != nil {
			return nil, err
		}
		matches = append(matches, m)
	}

	return matches, nil
}

func (s *SimulatorService) GetPointsUpToWeek(teamID, week int) (int, error) {
	var points int
	query := `
        SELECT COALESCE(SUM(
            CASE 
                WHEN (home_team_id = ? AND home_goals > away_goals) OR (away_team_id = ? AND away_goals > home_goals) THEN 3
                WHEN home_goals = away_goals THEN 1
                ELSE 0
            END), 0)
        FROM matches
        WHERE (home_team_id = ? OR away_team_id = ?) AND week < ?
    `
	err := s.DB.QueryRow(query, teamID, teamID, teamID, teamID, week).Scan(&points)
	if err != nil {
		return 0, err
	}
	return points, nil
}

func (s *SimulatorService) GetTotalPoints(teamID int) (int, error) {
	var points int
	query := `
        SELECT COALESCE(SUM(
            CASE 
                WHEN (home_team_id = ? AND home_goals > away_goals) OR (away_team_id = ? AND away_goals > home_goals) THEN 3
                WHEN home_goals = away_goals THEN 1
                ELSE 0
            END), 0)
        FROM matches
        WHERE home_team_id = ? OR away_team_id = ?
    `
	err := s.DB.QueryRow(query, teamID, teamID, teamID, teamID).Scan(&points)
	if err != nil {
		return 0, err
	}
	return points, nil
}

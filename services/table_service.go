package services

import (
	"insider-case/models"
	"sort"
)

type TeamStats struct {
	Team     models.Team
	Played   int
	Won      int
	Drawn    int
	Lost     int
	GF       int
	GA       int
	Points   int
	GoalDiff int
}

func (s *SimulatorService) GetCurrentStandings() ([]TeamStats, error) {
	rows, err := s.DB.Query(`
		SELECT home_team_id, away_team_id, home_goals, away_goals
		FROM matches
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	teams, err := s.getTeams()
	if err != nil {
		return nil, err
	}

	stats := make(map[int]*TeamStats)
	for _, t := range teams {
		stats[t.ID] = &TeamStats{Team: t}
	}

	for rows.Next() {
		var homeID, awayID, homeGoals, awayGoals int
		if err := rows.Scan(&homeID, &awayID, &homeGoals, &awayGoals); err != nil {
			return nil, err
		}

		homeStats := stats[homeID]
		awayStats := stats[awayID]

		homeStats.Played++
		awayStats.Played++
		homeStats.GF += homeGoals
		homeStats.GA += awayGoals
		awayStats.GF += awayGoals
		awayStats.GA += homeGoals

		homeStats.GoalDiff = homeStats.GF - homeStats.GA
		awayStats.GoalDiff = awayStats.GF - awayStats.GA

		switch {
		case homeGoals > awayGoals:
			homeStats.Won++
			homeStats.Points += 3
			awayStats.Lost++
		case homeGoals < awayGoals:
			awayStats.Won++
			awayStats.Points += 3
			homeStats.Lost++
		default:
			homeStats.Drawn++
			homeStats.Points++
			awayStats.Drawn++
			awayStats.Points++
		}
	}

	var standings []TeamStats
	for _, s := range stats {
		// Güncellenmiş istatistikleri Team struct'una da kopyala
		s.Team.Played = s.Played
		s.Team.Won = s.Won
		s.Team.Drawn = s.Drawn
		s.Team.Lost = s.Lost
		s.Team.Points = s.Points
		s.Team.GF = s.GF
		s.Team.GA = s.GA
		s.Team.GD = s.GoalDiff

		standings = append(standings, *s)
	}

	sort.SliceStable(standings, func(i, j int) bool {
		if standings[i].Points != standings[j].Points {
			return standings[i].Points > standings[j].Points
		}
		if standings[i].GoalDiff != standings[j].GoalDiff {
			return standings[i].GoalDiff > standings[j].GoalDiff
		}
		return standings[i].GF > standings[j].GF
	})

	return standings, nil
}

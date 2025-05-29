package router

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"insider-case/models"
	"insider-case/services"
)

type MatchHandler struct {
	simulator *services.SimulatorService
}

func NewMatchHandler(db *sql.DB) *MatchHandler {
	return &MatchHandler{
		simulator: services.NewSimulatorService(db),
	}
}

// POST /matches/simulate?week=1
// Belirtilen haftanın maçlarını simüle edip sonucu kaydeder ve sonucu döner
func (h *MatchHandler) SimulateMatchesHandler(w http.ResponseWriter, r *http.Request) {
	weekStr := r.URL.Query().Get("week")
	if weekStr == "" {
		http.Error(w, "Missing 'week' query parameter", http.StatusBadRequest)
		return
	}
	week, err := strconv.Atoi(weekStr)
	if err != nil || week < 1 {
		http.Error(w, "'week' must be a positive integer", http.StatusBadRequest)
		return
	}

	err = h.simulator.SimulateWeek(week)
	if err != nil {
		http.Error(w, "Failed to simulate matches: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Matches simulated for week " + weekStr))
}

// GET /matches?week=1
// İsteğe bağlı olarak haftaya göre maçları listeler, yoksa tüm maçları döner
func (h *MatchHandler) ListMatchesHandler(w http.ResponseWriter, r *http.Request) {
	weekStr := r.URL.Query().Get("week")

	var matches []models.Match
	var err error

	if weekStr == "" {
		matches, err = h.simulator.GetAllMatches()
	} else {
		week, err2 := strconv.Atoi(weekStr)
		if err2 != nil || week < 1 {
			http.Error(w, "'week' must be a positive integer", http.StatusBadRequest)
			return
		}
		matches, err = h.simulator.GetMatchesByWeek(week)
	}

	if err != nil {
		http.Error(w, "Failed to fetch matches: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matches)
}



package router

import (
	"database/sql"
	"encoding/json"
	"insider-case/services"
	"net/http"
)

type TableHandler struct {
	simulator *services.SimulatorService
}

func NewTableHandler(db *sql.DB) *TableHandler {
	return &TableHandler{
		simulator: services.NewSimulatorService(db),
	}
}

// GET /standings
// Güncel puan tablosunu JSON formatında döner
func (h *TableHandler) StandingsHandler(w http.ResponseWriter, r *http.Request) {
	standings, err := h.simulator.GetCurrentStandings()
	if err != nil {
		http.Error(w, "Failed to get standings: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(standings)
}

func ResetTable(db *sql.DB) error {
	// Tüm maçları sil
	_, err := db.Exec(`DELETE FROM matches`)
	if err != nil {
		return err
	}

	// Takım skorlarını sıfırla (id, name, strength hariç)
	_, err = db.Exec(`
		UPDATE teams SET 
			played = 0, 
			won = 0, 
			drawn = 0, 
			lost = 0, 
			points = 0,
			goal_difference = 0
	`)
	if err != nil {
		return err
	}

	return nil
}

package router

import (
	"database/sql"
	"encoding/json"
	"insider-case/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Router struct {
	db        *sql.DB
	simulator *services.SimulatorService
}

func NewRouter(db *sql.DB) *Router {
	return &Router{
		db:        db,
		simulator: services.NewSimulatorService(db),
	}
}

func (r *Router) SetupRoutes() http.Handler {
	mux := mux.NewRouter()

	mux.HandleFunc("/simulate/week", r.SimulateWeekHandler).Methods("POST")
	mux.HandleFunc("/simulate/all", r.SimulateAllHandler).Methods("POST")
	mux.HandleFunc("/standings", r.StandingsHandler).Methods("GET")

	return mux
}

// /simulate/week endpointi, istekte "week" parametresi bekliyor (POST form veya query param)
// İlgili haftayı simüle eder ve sonucu döner
func (r *Router) SimulateWeekHandler(w http.ResponseWriter, req *http.Request) {
	weekStr := req.URL.Query().Get("week")
	if weekStr == "" {
		http.Error(w, "Missing 'week' query parameter", http.StatusBadRequest)
		return
	}

	week, err := strconv.Atoi(weekStr)
	if err != nil || week < 1 {
		http.Error(w, "'week' must be a positive integer", http.StatusBadRequest)
		return
	}

	err = r.simulator.SimulateWeek(week)
	if err != nil {
		http.Error(w, "Failed to simulate week: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Week " + weekStr + " simulated successfully"))
}

// /simulate/all endpointi tüm haftaları simüle eder (örneğin 1-5 hafta)
func (r *Router) SimulateAllHandler(w http.ResponseWriter, req *http.Request) {
	err := r.simulator.SimulateAllWeeks()
	if err != nil {
		http.Error(w, "Failed to simulate all weeks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("All weeks simulated successfully"))
}

// /standings endpointi güncel puan tablosunu JSON formatında döner
func (r *Router) StandingsHandler(w http.ResponseWriter, req *http.Request) {
	standings, err := r.simulator.GetCurrentStandings()
	if err != nil {
		http.Error(w, "Failed to get standings: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(standings)
}

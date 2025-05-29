package models

type Match struct {
	ID         int    // Maç ID
	Week       int    // Haftası
	HomeTeamID int    // Ev sahibi takımın ID'si
	AwayTeamID int    // Deplasman takımının ID'si
	HomeGoals  int    // Ev sahibi takımın attığı gol sayısı
	AwayGoals  int    // Deplasman takımının attığı gol sayısı
	Result     string // "HomeWin", "AwayWin", "Draw" - Opsiyonel, hesaplanabilir
}

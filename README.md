# Insider-Backend-Hiringday-Task

A Go-based backend project demonstrating API development and SQLite database integration.  
Includes custom routing and HTTP server setup as part of the Insider Development Intern Hiring Day case study.

---

## 🚀 Introduction

This project implements a backend system to simulate a football league. The system fulfills the following requirements:

- Store team and match data in a structured format.  
- Simulate match results based on the relative strength of the teams.  
- Simulate matches on a **weekly** basis according to a schedule.  
- Provide an API to retrieve current league standings in real-time.  
- Include functionality to reset the system for a new simulation.

All weekly simulation logic is implemented using the Go programming language, simulating one match week at a time and updating team statistics accordingly.

---

## 🏗️ Solution Structure

The backend is organized into several packages to ensure clean separation of concerns:

- `handlers/` → HTTP handlers for API endpoints (`/match`, `/standings`)  
- `services/` → Business logic for simulations and league table calculations  
- `models/` → Data structures such as Team and Match  
- `router/` → Router setup and HTTP endpoint registration  
- `league.db` → SQLite database storing all league data  
- `reset.sql` → SQL script to reset the database for a new simulation  

---

## 🧠 Simulation Logic

- Each team has a **strength score**, derived from actual Premier League performance.  
- Matches are simulated **weekly**, with outcomes influenced by team strength.  
- League standings update automatically after every simulated week.

---

## 🏟️ Database Design

SQLite3 is used for data persistence. The database is managed and verified with **DB Browser for SQLite**.

### Tables

#### `matches` Table

Stores match details and weekly results between teams:

- `id`: Match unique identifier (auto-incremented)  
- `week`: The week number the match belongs to  
- `home_team_id`: The ID of the home team (foreign key referencing `teams.id`)  
- `away_team_id`: The ID of the away team (foreign key referencing `teams.id`)  
- `home_goals`: Goals scored by the home team (default 0)  
- `away_goals`: Goals scored by the away team (default 0)  
- `result`: Match result summary (e.g., "Home Win", "Draw", "Away Win")  

SQL for the `matches` table:

```sql
CREATE TABLE matches (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    week INTEGER NOT NULL,
    home_team_id INTEGER NOT NULL,
    away_team_id INTEGER NOT NULL,
    home_goals INTEGER DEFAULT 0,
    away_goals INTEGER DEFAULT 0,
    result TEXT,
    FOREIGN KEY(home_team_id) REFERENCES teams(id),
    FOREIGN KEY(away_team_id) REFERENCES teams(id)
);

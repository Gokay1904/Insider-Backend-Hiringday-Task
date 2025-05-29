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
.
├── handlers/
│   ├── match_handler.go       # Match simulation endpoint
│   └── table_handler.go       # Standings retrieval endpoint
├── models/
│   ├── interface.go           # Interface definitions
│   ├── match.go               # Match model
│   └── team.go                # Team model
├── router/
│   └── router.go              # HTTP router setup
├── services/
│   ├── match_service.go       # Business logic for matches
│   ├── simulation_service.go  # Match result calculation
│   └── table_service.go       # Standings and table logic
├── league.db                  # SQLite database file
├── reset.sql                  # SQL script to reset DB (truncate matches and standings)
├── main.go                    # Application entry point
├── go.mod                     # Go module file
├── go.sum                     # Go module dependencies
└── Dockerfile                 # Docker container configuration

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

#### `teams` Table

```sql
CREATE TABLE teams (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    position INTEGER DEFAULT 0,
    played INTEGER DEFAULT 0,
    won INTEGER DEFAULT 0,
    drawn INTEGER DEFAULT 0,
    lost INTEGER DEFAULT 0,
    gf INTEGER DEFAULT 0,
    ga INTEGER DEFAULT 0,
    gd INTEGER DEFAULT 0,
    points INTEGER DEFAULT 0,
    strength INTEGER NOT NULL
);

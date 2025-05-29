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

## 🗂️ Directory Structure

| Component           | Type       | Description                                                                 |
|---------------------|------------|-----------------------------------------------------------------------------|
| **`handlers/`**     | Package    | API endpoint controllers                                                    |
| **`models/`**       | Package    | Data models and domain interfaces                                           |
| **`router/`**       | Package    | HTTP routing configuration                                                  |
| **`services/`**     | Package    | Core business logic services                                                |
| **`database/`**     | Directory  | Database management files                                                   |
| **`main.go`**       | File       | Application entry point                                                     |
| **`Dockerfile`**    | File       | Container configuration                                                     |

## 📄 File Details


### 🖥️ Handlers Package

| File                | Responsibilities                                                                 |
|---------------------|---------------------------------------------------------------------------------|
| `match_handler.go`  | - Processes `/simulate/week` and `/simulate/all` requests<br>- Validates inputs<br>- Formats JSON responses |
| `table_handler.go`  | - Handles GET `/standings` requests<br>- Implements response caching<br>- Error handling |

### 🧠 Models Package

| File          | Key Components                                                                 |
|---------------|-------------------------------------------------------------------------------|
| `team.go`     | - `Team` struct<br>- `CalculatePoints()`<br>- `UpdateStats()`<br>- Validation logic |
| `match.go`    | - `Match` struct<br>- `Simulate()` method<br>- Result enums (HOME_WIN, etc.)      |
| `interface.go`| - `TeamRepository` interface<br>- `MatchRepository` interface                     |

### ⚙️ Services Package

| File                   | Core Logic                                                                 |
|------------------------|----------------------------------------------------------------------------|
| `match_service.go`     | - Week progression<br>- Team stats updates<br>- Simulation coordination    |
| `simulation_service.go`| - Probability algorithms<br>- Strength-based calculations<br>- Goal generation |
| `table_service.go`     | - Standings calculation<br>- Sorting logic<br>- Position assignment        |

### 🗃️ Database Files

| File        | Purpose                                                                 |
|-------------|-------------------------------------------------------------------------|
| `schema.sql`| - Table creation<br>- Indexes and constraints<br>- Database versioning  |
| `seed.sql`  | - Initial team data<br>- Match schedule<br>- Realistic strength ratings |
| `reset.sql` | - Data truncation<br>- Statistics reset<br>- Referential integrity      |
---

## 🧠 Simulation Logic

- Each team has a **strength score**, derived from actual Premier League performance.  
- Matches are simulated **weekly**, with outcomes influenced by team strength.

### ➤ Normalize Team Strength
Each team's strength is scaled between 0 and 1:

- `homeFactor = homeStrength / 100`
- `awayFactor = awayStrength / 100`

### ➤ Calculate Expected Goals
Expected goals (λ) are computed using team strength and home advantage:

- `homeLambda = 1.8 * homeFactor`
- `awayLambda = 1.0 * awayFactor`

### ➤ Simulate Goals with Poisson Distribution
The number of goals is sampled using the Poisson distribution:

- `P(k; λ) = (e^-λ * λ^k) / k!`

Where:
- `P(k; λ)` is the probability of scoring exactly `k` goals
- `λ` is the expected goal count for that team

### ➤ Score Limiting
To ensure realism, the number of goals is capped at **5 per team**.

- League standings update automatically after every simulated week.

---

## 🏟️ Database Design

SQLite3 is used for data persistence. The database is managed and verified with **DB Browser for SQLite**.

### Tables

| Table Name | Columns                                                                                                                | Description                         |
|------------|------------------------------------------------------------------------------------------------------------------------|-----------------------------------|
| **teams**  | `id` (INTEGER, PK, AUTOINCREMENT), `name` (TEXT, NOT NULL), `position` (INTEGER), `played` (INTEGER), `won` (INTEGER), `drawn` (INTEGER), `lost` (INTEGER), `gf` (INTEGER), `ga` (INTEGER), `gd` (INTEGER), `points` (INTEGER), `strength` (INTEGER, NOT NULL) | Stores team info and stats         |
| **matches**| `id` (INTEGER, PK, AUTOINCREMENT), `week` (INTEGER, NOT NULL), `home_team_id` (INTEGER, FK), `away_team_id` (INTEGER, FK), `home_goals` (INTEGER), `away_goals` (INTEGER), `result` (TEXT) | Stores match info and results      |

---
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
```

#### `matches` Table
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
```


#### `seed` code for initial conditions of selected teams
```sql

-- Clear existing data
DELETE FROM teams;
DELETE FROM matches;

-- Insert 5 teams with different strengths (example data)
INSERT INTO teams (name, position, played, won, drawn, lost, gf, ga, gd, points, strength) VALUES
('Arsenal', 1, 0, 0, 0, 0, 0, 0, 0, 0, 85),
('Manchester City', 2, 0, 0, 0, 0, 0, 0, 0, 0, 90),
('Manchester United', 3, 0, 0, 0, 0, 0, 0, 0, 0, 78),
('Chelsea', 4, 0, 0, 0, 0, 0, 0, 0, 0, 80),
('Liverpool', 5, 0, 0, 0, 0, 0, 0, 0, 0, 88);
```

## 🚀 Available Endpoints

| Endpoint         | Method | Description                   | Request Body | Response                    |
|------------------|--------|------------------------------|--------------|-----------------------------|
| `/simulate/week` | POST   | Simulates next week's matches | None         | JSON: Simulated matches     |
| `/simulate/all`  | POST   | Simulates all remaining weeks | None         | JSON: All simulated matches |
| `/standings`     | GET    | Returns current league table  | None         | JSON: Team standings        |
| `/reset`         | POST   | Resets all matches and stats  | None         | Plain text confirmation     |

---

📝 Note: All match results and league updates triggered by these endpoints are automatically persisted in the SQLite database (league.db) under the teams and matches tables.

### How to Call Endpoints with `curl`

- **Simulate a specific week**

To play weekly
  ```bash
  curl -X POST "http://localhost:8080/simulate/week?week=2"
  ```

To simulate all
   ```bash
  curl -X POST "http://localhost:8080/simulate/all"
  ```

To reset matches

 ```bash
  curl -X POST "http://localhost:8080/reset
  ```


## ⚙️ Setup and Running the Project (Without Docker)

Follow these steps to get the project up and running on your local machine:

### Prerequisites

- **Go** installed (version 1.18 or higher recommended)  
  Download from: https://golang.org/dl/  
- **SQLite3** installed (for database management)  
  Download from: https://sqlite.org/download.html  
- (Optional) A SQLite database browser like **DB Browser for SQLite** for inspecting and managing the database visually  
  https://sqlitebrowser.org/

---

### Step 1: Clone the repository

```bash
git clone <repository-url>
cd Insider-Backend-Hiringday-Task
```

### Step 2: Install Go dependencies

Run the following command in the project root to download all necessary Go modules:

```bash
go mod tidy
```

### Step 3: Prepare the SQLite database

Make sure `sqlite3` is installed on your system.

Run the following command to initialize or reset the database using the provided SQL script:

```bash
sqlite3 league.db < reset.sql
```

Create tables under league database with teams and matches:
```bash
sqlite3 league.db < schema.sql
```

At the start of the simulation, each team—such as Arsenal—is initialized with predefined attributes including starting scores, match history (if applicable), and a strength rating. This seed data represents the initial conditions for each team and can be easily adjusted to simulate different league dynamics or difficulty levels.
```bash
sqlite3 league.db < seed.sql
```

### Step 4: Run the Backend Server

Start the backend application by running the following command in the project root directory:

```bash
go run main.go
```
This will start the server which will be listen on:
http://localhost:8080


### Step 5: Test the API Endpoints

You can test the backend API using `curl` commands or any API client like Postman or Insomnia.

To play weekly (for example week = 2)
  ```bash
  curl -X POST "http://localhost:8080/simulate/week?week=2"
  ```

To simulate all
   ```bash
  curl -X POST "http://localhost:8080/simulate/all"
  ```

To manually reset matches

 ```bash
  curl -X POST "http://localhost:8080/reset
  ```


### Output examples:

First predictions:
```yaml
Week 1 - Championship Predictions
Arsenal: 30.00
Manchester United: 25.00
Manchester City: 20.00
Chelsea: 15.00
Liverpool: 10.00

   ```

Then the match results:
   ```yaml
Week 1 - Match Results
- match:
    home_team: Arsenal
    away_team: Chelsea
    home_goals: 2
    away_goals: 1
- match:
    home_team: Liverpool
    away_team: Manchester United
    home_goals: 0
    away_goals: 3
- bye: Manchester City

   ```
Finally the team scores after the match is played:

  ```yaml
teams:
  - name: Team A
    stats:
      Played: 1
      Won: 1
      Drawn: 0
      Lost: 0
      GF: 3
      GA: 1
      GoalDiff: 2
      Points: 3

  - name: Team B
    stats:
      Played: 1
      Won: 0
      Drawn: 0
      Lost: 1
      GF: 1
      GA: 3
      GoalDiff: -2
      Points: 0

 ```

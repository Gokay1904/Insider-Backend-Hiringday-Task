-- Eğer varsa eski tabloları sil
DROP TABLE IF EXISTS matches;
DROP TABLE IF EXISTS teams;

-- Takımlar tablosu
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

-- Maçlar tablosu
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

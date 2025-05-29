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

-- Tüm maç kayıtlarını sil
DELETE FROM matches;

-- Tüm takım istatistiklerini sıfırla (id, name, strength hariç)
UPDATE teams
SET
    position = 0,
    played = 0,
    won = 0,
    drawn = 0,
    lost = 0,
    gf = 0,
    ga = 0,
    gd = 0,
    points = 0;

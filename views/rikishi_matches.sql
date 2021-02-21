CREATE VIEW rikishi_matches AS
SELECT 
	r.id as rikishi_id,
	m.id as match_id,
	t.id as tournament_id,
	m.day as day,
    m.winner_id != 0 as concluded,
	m.winner_id = r.id as won,
	r2.id as opponent
FROM matches as m
JOIN tournaments as t
	ON m.tournament_id = t.id
JOIN rikishis as r
	ON m.east_id = r.id
JOIN rikishis as r2
	ON m.west_id = r2.id
UNION
SELECT
	r.id as rikishi_id,
	m.id as match_id,
	t.id as tournament_id,
	m.day as day,
    m.winner_id != 0 as concluded,
	m.winner_id = r.id as won,
	r2.id as opponent
FROM matches as m
JOIN tournaments as t
	ON m.tournament_id = t.id
JOIN rikishis as r
	ON m.west_id = r.id
JOIN rikishis as r2
	ON m.east_id = r2.id
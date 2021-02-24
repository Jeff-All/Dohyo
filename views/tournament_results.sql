CREATE VIEW tournament_results AS
SELECT
	t.name as tournament,
	r.id as rikishi_id,
	rm.tournament_id,
	sum(rm.won) as wins,
	sum(CASE WHEN rm.won = 0 THEN 1 ELSE 0 END) as losses
FROM rikishis as r
JOIN rikishi_matches as rm
	ON r.id = rm.rikishi_id
	AND rm.concluded = 1
JOIN tournaments as t
	ON t.id = rm.tournament_id
GROUP BY r.id, rm.tournament_id
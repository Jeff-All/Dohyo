CREATE VIEW tournament_results AS
SELECT
	r.id,
	rm.tournament_id,
	sum(rm.won) as wins,
	sum(CASE WHEN rm.won = 0 THEN 1 ELSE 0 END) as losses
FROM rikishis as r
JOIN rikishi_matches as rm
	ON r.id = rm.rikishi_id
	AND rm.concluded = 1
GROUP BY r.id, rm.tournament_id
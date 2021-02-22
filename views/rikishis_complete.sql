CREATE VIEW rikishis_complete as
SELECT
	r.id,
	r.name,
	r.avatar,
	CASE WHEN r.sub_rank > 0 THEN ranks.name || " " || r.sub_rank ELSE ranks.name END as rank
FROM rikishis as r
JOIN ranks
ON ranks.id = r.rank_id
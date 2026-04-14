-- Simplify matches to single-game records and track total points per user.

DROP TABLE IF EXISTS match_sets;

-- Existing match rows can't be mapped to the new single-score schema, so drop them.
DELETE FROM matches;

ALTER TABLE matches DROP COLUMN IF EXISTS match_format;
ALTER TABLE matches DROP COLUMN IF EXISTS player1_sets_won;
ALTER TABLE matches DROP COLUMN IF EXISTS player2_sets_won;
ALTER TABLE matches DROP COLUMN IF EXISTS tournament_match_id;

ALTER TABLE matches ADD COLUMN player1_score SMALLINT NOT NULL DEFAULT 0 CHECK (player1_score >= 0);
ALTER TABLE matches ADD COLUMN player2_score SMALLINT NOT NULL DEFAULT 0 CHECK (player2_score >= 0);
ALTER TABLE matches ADD CONSTRAINT different_scores CHECK (player1_score <> player2_score);

-- Reset aggregates because match history was just cleared.
UPDATE users SET matches_played = 0, matches_won = 0;

ALTER TABLE users ADD COLUMN total_points INTEGER NOT NULL DEFAULT 0 CHECK (total_points >= 0);

ALTER TABLE users DROP COLUMN IF EXISTS total_points;

ALTER TABLE matches DROP CONSTRAINT IF EXISTS different_scores;
ALTER TABLE matches DROP COLUMN IF EXISTS player1_score;
ALTER TABLE matches DROP COLUMN IF EXISTS player2_score;

ALTER TABLE matches ADD COLUMN match_format SMALLINT NOT NULL DEFAULT 5 CHECK (match_format IN (3, 5, 7));
ALTER TABLE matches ADD COLUMN player1_sets_won SMALLINT NOT NULL DEFAULT 0;
ALTER TABLE matches ADD COLUMN player2_sets_won SMALLINT NOT NULL DEFAULT 0;
ALTER TABLE matches ADD COLUMN tournament_match_id UUID;

CREATE TABLE match_sets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    set_number SMALLINT NOT NULL CHECK (set_number BETWEEN 1 AND 7),
    player1_score SMALLINT NOT NULL CHECK (player1_score >= 0),
    player2_score SMALLINT NOT NULL CHECK (player2_score >= 0),
    UNIQUE (match_id, set_number)
);

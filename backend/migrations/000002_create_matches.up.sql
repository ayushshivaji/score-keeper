CREATE TABLE matches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player1_id UUID NOT NULL REFERENCES users(id),
    player2_id UUID NOT NULL REFERENCES users(id),
    winner_id UUID NOT NULL REFERENCES users(id),
    match_format SMALLINT NOT NULL CHECK (match_format IN (3, 5, 7)),
    player1_sets_won SMALLINT NOT NULL,
    player2_sets_won SMALLINT NOT NULL,
    tournament_match_id UUID,
    played_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL REFERENCES users(id),
    CONSTRAINT different_players CHECK (player1_id <> player2_id)
);

CREATE TABLE match_sets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    set_number SMALLINT NOT NULL CHECK (set_number BETWEEN 1 AND 7),
    player1_score SMALLINT NOT NULL CHECK (player1_score >= 0),
    player2_score SMALLINT NOT NULL CHECK (player2_score >= 0),
    UNIQUE (match_id, set_number)
);

CREATE INDEX idx_matches_player1 ON matches(player1_id);
CREATE INDEX idx_matches_player2 ON matches(player2_id);
CREATE INDEX idx_matches_played_at ON matches(played_at DESC);

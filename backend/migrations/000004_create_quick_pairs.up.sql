CREATE TABLE quick_pairs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    player1_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    player2_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (player1_id <> player2_id)
);

-- LEAST/GREATEST keeps (A,B) and (B,A) from both existing for the same user.
CREATE UNIQUE INDEX quick_pairs_user_pair_unique
ON quick_pairs (user_id, LEAST(player1_id, player2_id), GREATEST(player1_id, player2_id));

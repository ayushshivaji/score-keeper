package repository

import (
	"context"
	"time"

	"github.com/ayush-sr/score-keeper/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

const userCols = `id, google_id, email, name, avatar_url, matches_played, matches_won, total_points, created_at, updated_at`

func scanUser(row pgx.Row, u *model.User) error {
	return row.Scan(
		&u.ID, &u.GoogleID, &u.Email, &u.Name, &u.AvatarURL,
		&u.MatchesPlayed, &u.MatchesWon, &u.TotalPoints, &u.CreatedAt, &u.UpdatedAt,
	)
}

func (r *UserRepository) UpsertByGoogleID(ctx context.Context, googleID, email, name string, avatarURL *string) (*model.User, error) {
	var user model.User
	err := scanUser(r.db.QueryRow(ctx, `
		INSERT INTO users (google_id, email, name, avatar_url)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (google_id) DO UPDATE SET
			email = EXCLUDED.email,
			name = EXCLUDED.name,
			avatar_url = EXCLUDED.avatar_url,
			updated_at = NOW()
		RETURNING `+userCols+`
	`, googleID, email, name, avatarURL), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := scanUser(r.db.QueryRow(ctx, `SELECT `+userCols+` FROM users WHERE id = $1`, id), &user)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) List(ctx context.Context, search string, page, perPage int) ([]model.User, int, error) {
	offset := (page - 1) * perPage

	var total int
	var rows pgx.Rows
	var err error

	if search != "" {
		pattern := "%" + search + "%"
		if err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE name ILIKE $1`, pattern).Scan(&total); err != nil {
			return nil, 0, err
		}
		rows, err = r.db.Query(ctx, `SELECT `+userCols+` FROM users WHERE name ILIKE $1 ORDER BY name LIMIT $2 OFFSET $3`,
			pattern, perPage, offset)
	} else {
		if err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total); err != nil {
			return nil, 0, err
		}
		rows, err = r.db.Query(ctx, `SELECT `+userCols+` FROM users ORDER BY name LIMIT $1 OFFSET $2`, perPage, offset)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := scanUser(rows, &u); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, nil
}

func (r *UserRepository) ListLeaderboard(ctx context.Context, page, perPage int) ([]model.User, int, error) {
	offset := (page - 1) * perPage

	var total int
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE matches_played > 0`).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT `+userCols+`
		FROM users
		WHERE matches_played > 0
		ORDER BY matches_won DESC,
		         total_points DESC,
		         (matches_won::float / NULLIF(matches_played, 0)) DESC NULLS LAST,
		         name ASC
		LIMIT $1 OFFSET $2
	`, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := scanUser(rows, &u); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, nil
}

func (r *UserRepository) StoreRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`, userID, tokenHash, expiresAt)
	return err
}

func (r *UserRepository) GetRefreshToken(ctx context.Context, tokenHash string) (uuid.UUID, time.Time, error) {
	var userID uuid.UUID
	var expiresAt time.Time
	err := r.db.QueryRow(ctx, `
		SELECT user_id, expires_at FROM refresh_tokens WHERE token_hash = $1
	`, tokenHash).Scan(&userID, &expiresAt)
	if err == pgx.ErrNoRows {
		return uuid.Nil, time.Time{}, nil
	}
	return userID, expiresAt, err
}

func (r *UserRepository) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM refresh_tokens WHERE token_hash = $1`, tokenHash)
	return err
}

func (r *UserRepository) DeleteUserRefreshTokens(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM refresh_tokens WHERE user_id = $1`, userID)
	return err
}

// AdjustMatchStats applies signed deltas to a player's aggregate counters.
// Use +1 / positive values when recording a match, -1 / negative on deletion.
func (r *UserRepository) AdjustMatchStats(ctx context.Context, tx pgx.Tx, playerID uuid.UUID, playedDelta, wonDelta, pointsDelta int) error {
	_, err := tx.Exec(ctx, `
		UPDATE users
		SET matches_played = matches_played + $2,
		    matches_won = matches_won + $3,
		    total_points = total_points + $4,
		    updated_at = NOW()
		WHERE id = $1
	`, playerID, playedDelta, wonDelta, pointsDelta)
	return err
}

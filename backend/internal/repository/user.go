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

func (r *UserRepository) UpsertByGoogleID(ctx context.Context, googleID, email, name string, avatarURL *string) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(ctx, `
		INSERT INTO users (google_id, email, name, avatar_url)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (google_id) DO UPDATE SET
			email = EXCLUDED.email,
			name = EXCLUDED.name,
			avatar_url = EXCLUDED.avatar_url,
			updated_at = NOW()
		RETURNING id, google_id, email, name, avatar_url, matches_played, matches_won, created_at, updated_at
	`, googleID, email, name, avatarURL).Scan(
		&user.ID, &user.GoogleID, &user.Email, &user.Name, &user.AvatarURL,
		&user.MatchesPlayed, &user.MatchesWon, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(ctx, `
		SELECT id, google_id, email, name, avatar_url, matches_played, matches_won, created_at, updated_at
		FROM users WHERE id = $1
	`, id).Scan(
		&user.ID, &user.GoogleID, &user.Email, &user.Name, &user.AvatarURL,
		&user.MatchesPlayed, &user.MatchesWon, &user.CreatedAt, &user.UpdatedAt,
	)
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
		err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE name ILIKE $1`, pattern).Scan(&total)
		if err != nil {
			return nil, 0, err
		}
		rows, err = r.db.Query(ctx, `
			SELECT id, google_id, email, name, avatar_url, matches_played, matches_won, created_at, updated_at
			FROM users WHERE name ILIKE $1 ORDER BY name LIMIT $2 OFFSET $3
		`, pattern, perPage, offset)
	} else {
		err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total)
		if err != nil {
			return nil, 0, err
		}
		rows, err = r.db.Query(ctx, `
			SELECT id, google_id, email, name, avatar_url, matches_played, matches_won, created_at, updated_at
			FROM users ORDER BY name LIMIT $1 OFFSET $2
		`, perPage, offset)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.GoogleID, &u.Email, &u.Name, &u.AvatarURL,
			&u.MatchesPlayed, &u.MatchesWon, &u.CreatedAt, &u.UpdatedAt); err != nil {
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

func (r *UserRepository) IncrementMatchStats(ctx context.Context, tx pgx.Tx, playerID uuid.UUID, won bool) error {
	if won {
		_, err := tx.Exec(ctx, `
			UPDATE users SET matches_played = matches_played + 1, matches_won = matches_won + 1, updated_at = NOW()
			WHERE id = $1
		`, playerID)
		return err
	}
	_, err := tx.Exec(ctx, `
		UPDATE users SET matches_played = matches_played + 1, updated_at = NOW()
		WHERE id = $1
	`, playerID)
	return err
}

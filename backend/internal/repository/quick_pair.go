package repository

import (
	"context"
	"errors"

	"github.com/ayush-sr/score-keeper/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QuickPairRepository struct {
	db *pgxpool.Pool
}

func NewQuickPairRepository(db *pgxpool.Pool) *QuickPairRepository {
	return &QuickPairRepository{db: db}
}

// listSelect joins twice onto users so each row carries both player structs.
// userCols is imported (package-local) from user.go.
const quickPairCols = `qp.id, qp.user_id, qp.player1_id, qp.player2_id, qp.created_at,
	p1.id, p1.google_id, p1.email, p1.name, p1.avatar_url, p1.matches_played, p1.matches_won, p1.total_points, p1.created_at, p1.updated_at,
	p2.id, p2.google_id, p2.email, p2.name, p2.avatar_url, p2.matches_played, p2.matches_won, p2.total_points, p2.created_at, p2.updated_at`

func scanQuickPairWithPlayers(row pgx.Row, qp *model.QuickPairWithPlayers) error {
	return row.Scan(
		&qp.ID, &qp.UserID, &qp.Player1ID, &qp.Player2ID, &qp.CreatedAt,
		&qp.Player1.ID, &qp.Player1.GoogleID, &qp.Player1.Email, &qp.Player1.Name, &qp.Player1.AvatarURL, &qp.Player1.MatchesPlayed, &qp.Player1.MatchesWon, &qp.Player1.TotalPoints, &qp.Player1.CreatedAt, &qp.Player1.UpdatedAt,
		&qp.Player2.ID, &qp.Player2.GoogleID, &qp.Player2.Email, &qp.Player2.Name, &qp.Player2.AvatarURL, &qp.Player2.MatchesPlayed, &qp.Player2.MatchesWon, &qp.Player2.TotalPoints, &qp.Player2.CreatedAt, &qp.Player2.UpdatedAt,
	)
}

func (r *QuickPairRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]model.QuickPairWithPlayers, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+quickPairCols+`
		FROM quick_pairs qp
		JOIN users p1 ON qp.player1_id = p1.id
		JOIN users p2 ON qp.player2_id = p2.id
		WHERE qp.user_id = $1
		ORDER BY qp.created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pairs := []model.QuickPairWithPlayers{}
	for rows.Next() {
		var qp model.QuickPairWithPlayers
		if err := scanQuickPairWithPlayers(rows, &qp); err != nil {
			return nil, err
		}
		pairs = append(pairs, qp)
	}
	return pairs, nil
}

// ErrQuickPairDuplicate is returned when the unique index on
// (user_id, pair) is violated by a Create.
var ErrQuickPairDuplicate = errors.New("quick pair already exists for this user")

func (r *QuickPairRepository) Create(ctx context.Context, userID, player1ID, player2ID uuid.UUID) (*model.QuickPairWithPlayers, error) {
	var id uuid.UUID
	err := r.db.QueryRow(ctx, `
		INSERT INTO quick_pairs (user_id, player1_id, player2_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`, userID, player1ID, player2ID).Scan(&id)
	if err != nil {
		// pgx v5 exposes SQLState() via a PgError; check for unique violation (23505).
		type pgErr interface{ SQLState() string }
		var pe pgErr
		if errors.As(err, &pe) && pe.SQLState() == "23505" {
			return nil, ErrQuickPairDuplicate
		}
		return nil, err
	}

	var qp model.QuickPairWithPlayers
	err = scanQuickPairWithPlayers(r.db.QueryRow(ctx, `
		SELECT `+quickPairCols+`
		FROM quick_pairs qp
		JOIN users p1 ON qp.player1_id = p1.id
		JOIN users p2 ON qp.player2_id = p2.id
		WHERE qp.id = $1
	`, id), &qp)
	if err != nil {
		return nil, err
	}
	return &qp, nil
}

func (r *QuickPairRepository) Delete(ctx context.Context, userID, id uuid.UUID) (bool, error) {
	tag, err := r.db.Exec(ctx, `DELETE FROM quick_pairs WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

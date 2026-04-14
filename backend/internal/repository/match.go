package repository

import (
	"context"
	"fmt"

	"github.com/ayush-sr/score-keeper/backend/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MatchRepository struct {
	db *pgxpool.Pool
}

func NewMatchRepository(db *pgxpool.Pool) *MatchRepository {
	return &MatchRepository{db: db}
}

func (r *MatchRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.db.Begin(ctx)
}

func (r *MatchRepository) CreateMatch(ctx context.Context, tx pgx.Tx, m *model.Match) error {
	return tx.QueryRow(ctx, `
		INSERT INTO matches (player1_id, player2_id, winner_id, player1_score, player2_score, played_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`, m.Player1ID, m.Player2ID, m.WinnerID, m.Player1Score, m.Player2Score,
		m.PlayedAt, m.CreatedBy).Scan(&m.ID, &m.CreatedAt)
}

const matchWithDetailsCols = `m.id, m.player1_id, m.player2_id, m.winner_id,
		m.player1_score, m.player2_score, m.played_at, m.created_at, m.created_by,
		p1.id, p1.google_id, p1.email, p1.name, p1.avatar_url, p1.matches_played, p1.matches_won, p1.total_points, p1.created_at, p1.updated_at,
		p2.id, p2.google_id, p2.email, p2.name, p2.avatar_url, p2.matches_played, p2.matches_won, p2.total_points, p2.created_at, p2.updated_at`

func scanMatchWithDetails(row pgx.Row, m *model.MatchWithDetails) error {
	return row.Scan(
		&m.ID, &m.Player1ID, &m.Player2ID, &m.WinnerID,
		&m.Player1Score, &m.Player2Score, &m.PlayedAt, &m.CreatedAt, &m.CreatedBy,
		&m.Player1.ID, &m.Player1.GoogleID, &m.Player1.Email, &m.Player1.Name, &m.Player1.AvatarURL, &m.Player1.MatchesPlayed, &m.Player1.MatchesWon, &m.Player1.TotalPoints, &m.Player1.CreatedAt, &m.Player1.UpdatedAt,
		&m.Player2.ID, &m.Player2.GoogleID, &m.Player2.Email, &m.Player2.Name, &m.Player2.AvatarURL, &m.Player2.MatchesPlayed, &m.Player2.MatchesWon, &m.Player2.TotalPoints, &m.Player2.CreatedAt, &m.Player2.UpdatedAt,
	)
}

func (r *MatchRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.MatchWithDetails, error) {
	var m model.MatchWithDetails
	query := `
		SELECT ` + matchWithDetailsCols + `
		FROM matches m
		JOIN users p1 ON m.player1_id = p1.id
		JOIN users p2 ON m.player2_id = p2.id
		WHERE m.id = $1
	`
	err := scanMatchWithDetails(r.db.QueryRow(ctx, query, id), &m)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MatchRepository) List(ctx context.Context, playerID *uuid.UUID, page, perPage int) ([]model.MatchWithDetails, int, error) {
	offset := (page - 1) * perPage

	baseQuery := `
		FROM matches m
		JOIN users p1 ON m.player1_id = p1.id
		JOIN users p2 ON m.player2_id = p2.id
	`
	where := ""
	args := []interface{}{}
	argIdx := 1

	if playerID != nil {
		where = fmt.Sprintf(" WHERE m.player1_id = $%d OR m.player2_id = $%d", argIdx, argIdx)
		args = append(args, *playerID)
		argIdx++
	}

	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)
	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) "+baseQuery+where, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, perPage, offset)
	query := fmt.Sprintf("SELECT %s %s%s ORDER BY m.played_at DESC LIMIT $%d OFFSET $%d",
		matchWithDetailsCols, baseQuery, where, argIdx, argIdx+1)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var matches []model.MatchWithDetails
	for rows.Next() {
		var m model.MatchWithDetails
		if err := scanMatchWithDetails(rows, &m); err != nil {
			return nil, 0, err
		}
		matches = append(matches, m)
	}
	return matches, total, nil
}

func (r *MatchRepository) Delete(ctx context.Context, tx pgx.Tx, id uuid.UUID) error {
	_, err := tx.Exec(ctx, `DELETE FROM matches WHERE id = $1`, id)
	return err
}

func (r *MatchRepository) GetCreatedBy(ctx context.Context, matchID uuid.UUID) (uuid.UUID, error) {
	var createdBy uuid.UUID
	err := r.db.QueryRow(ctx, `SELECT created_by FROM matches WHERE id = $1`, matchID).Scan(&createdBy)
	return createdBy, err
}

// GetPlayerResultsChronological returns the player's match outcomes (true = win)
// in order from oldest to newest. Used for streak and recent-form calculations.
func (r *MatchRepository) GetPlayerResultsChronological(ctx context.Context, playerID uuid.UUID) ([]bool, error) {
	rows, err := r.db.Query(ctx, `
		SELECT winner_id = $1 AS won
		FROM matches
		WHERE player1_id = $1 OR player2_id = $1
		ORDER BY played_at ASC, created_at ASC
	`, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []bool
	for rows.Next() {
		var won bool
		if err := rows.Scan(&won); err != nil {
			return nil, err
		}
		results = append(results, won)
	}
	return results, nil
}

// HeadToHeadAggregates returns totals for two players' shared match history.
func (r *MatchRepository) HeadToHeadAggregates(ctx context.Context, p1, p2 uuid.UUID) (total, p1Wins, p2Wins, p1Points, p2Points int, err error) {
	err = r.db.QueryRow(ctx, `
		SELECT
			COUNT(*),
			COALESCE(SUM(CASE WHEN winner_id = $1 THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN winner_id = $2 THEN 1 ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN player1_id = $1 THEN player1_score ELSE player2_score END), 0),
			COALESCE(SUM(CASE WHEN player1_id = $2 THEN player1_score ELSE player2_score END), 0)
		FROM matches
		WHERE (player1_id = $1 AND player2_id = $2)
		   OR (player1_id = $2 AND player2_id = $1)
	`, p1, p2).Scan(&total, &p1Wins, &p2Wins, &p1Points, &p2Points)
	return
}

// HeadToHeadMatches returns all matches between two players, newest first.
func (r *MatchRepository) HeadToHeadMatches(ctx context.Context, p1, p2 uuid.UUID) ([]model.MatchWithDetails, error) {
	query := `
		SELECT ` + matchWithDetailsCols + `
		FROM matches m
		JOIN users p1 ON m.player1_id = p1.id
		JOIN users p2 ON m.player2_id = p2.id
		WHERE (m.player1_id = $1 AND m.player2_id = $2)
		   OR (m.player1_id = $2 AND m.player2_id = $1)
		ORDER BY m.played_at DESC
	`
	rows, err := r.db.Query(ctx, query, p1, p2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []model.MatchWithDetails
	for rows.Next() {
		var m model.MatchWithDetails
		if err := scanMatchWithDetails(rows, &m); err != nil {
			return nil, err
		}
		matches = append(matches, m)
	}
	return matches, nil
}

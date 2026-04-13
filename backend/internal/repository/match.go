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
		INSERT INTO matches (player1_id, player2_id, winner_id, match_format, player1_sets_won, player2_sets_won, tournament_match_id, played_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at
	`, m.Player1ID, m.Player2ID, m.WinnerID, m.MatchFormat, m.Player1SetsWon, m.Player2SetsWon,
		m.TournamentMatchID, m.PlayedAt, m.CreatedBy).Scan(&m.ID, &m.CreatedAt)
}

func (r *MatchRepository) CreateMatchSets(ctx context.Context, tx pgx.Tx, matchID uuid.UUID, sets []model.MatchSet) error {
	for i := range sets {
		sets[i].MatchID = matchID
		sets[i].SetNumber = i + 1
		err := tx.QueryRow(ctx, `
			INSERT INTO match_sets (match_id, set_number, player1_score, player2_score)
			VALUES ($1, $2, $3, $4) RETURNING id
		`, matchID, sets[i].SetNumber, sets[i].Player1Score, sets[i].Player2Score).Scan(&sets[i].ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *MatchRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.MatchWithDetails, error) {
	var m model.MatchWithDetails
	err := r.db.QueryRow(ctx, `
		SELECT m.id, m.player1_id, m.player2_id, m.winner_id, m.match_format,
			m.player1_sets_won, m.player2_sets_won, m.tournament_match_id,
			m.played_at, m.created_at, m.created_by,
			p1.id, p1.google_id, p1.email, p1.name, p1.avatar_url, p1.matches_played, p1.matches_won, p1.created_at, p1.updated_at,
			p2.id, p2.google_id, p2.email, p2.name, p2.avatar_url, p2.matches_played, p2.matches_won, p2.created_at, p2.updated_at
		FROM matches m
		JOIN users p1 ON m.player1_id = p1.id
		JOIN users p2 ON m.player2_id = p2.id
		WHERE m.id = $1
	`, id).Scan(
		&m.ID, &m.Player1ID, &m.Player2ID, &m.WinnerID, &m.MatchFormat,
		&m.Player1SetsWon, &m.Player2SetsWon, &m.TournamentMatchID,
		&m.PlayedAt, &m.CreatedAt, &m.CreatedBy,
		&m.Player1.ID, &m.Player1.GoogleID, &m.Player1.Email, &m.Player1.Name, &m.Player1.AvatarURL, &m.Player1.MatchesPlayed, &m.Player1.MatchesWon, &m.Player1.CreatedAt, &m.Player1.UpdatedAt,
		&m.Player2.ID, &m.Player2.GoogleID, &m.Player2.Email, &m.Player2.Name, &m.Player2.AvatarURL, &m.Player2.MatchesPlayed, &m.Player2.MatchesWon, &m.Player2.CreatedAt, &m.Player2.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	sets, err := r.GetMatchSets(ctx, id)
	if err != nil {
		return nil, err
	}
	m.Sets = sets
	return &m, nil
}

func (r *MatchRepository) GetMatchSets(ctx context.Context, matchID uuid.UUID) ([]model.MatchSet, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, match_id, set_number, player1_score, player2_score
		FROM match_sets WHERE match_id = $1 ORDER BY set_number
	`, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sets []model.MatchSet
	for rows.Next() {
		var s model.MatchSet
		if err := rows.Scan(&s.ID, &s.MatchID, &s.SetNumber, &s.Player1Score, &s.Player2Score); err != nil {
			return nil, err
		}
		sets = append(sets, s)
	}
	return sets, nil
}

func (r *MatchRepository) List(ctx context.Context, playerID *uuid.UUID, page, perPage int) ([]model.MatchWithDetails, int, error) {
	offset := (page - 1) * perPage
	var total int
	var rows pgx.Rows
	var err error

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
	err = r.db.QueryRow(ctx, "SELECT COUNT(*) "+baseQuery+where, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	selectCols := `m.id, m.player1_id, m.player2_id, m.winner_id, m.match_format,
		m.player1_sets_won, m.player2_sets_won, m.tournament_match_id,
		m.played_at, m.created_at, m.created_by,
		p1.id, p1.google_id, p1.email, p1.name, p1.avatar_url, p1.matches_played, p1.matches_won, p1.created_at, p1.updated_at,
		p2.id, p2.google_id, p2.email, p2.name, p2.avatar_url, p2.matches_played, p2.matches_won, p2.created_at, p2.updated_at`

	args = append(args, perPage, offset)
	query := fmt.Sprintf("SELECT %s %s%s ORDER BY m.played_at DESC LIMIT $%d OFFSET $%d",
		selectCols, baseQuery, where, argIdx, argIdx+1)

	rows, err = r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var matches []model.MatchWithDetails
	for rows.Next() {
		var m model.MatchWithDetails
		if err := rows.Scan(
			&m.ID, &m.Player1ID, &m.Player2ID, &m.WinnerID, &m.MatchFormat,
			&m.Player1SetsWon, &m.Player2SetsWon, &m.TournamentMatchID,
			&m.PlayedAt, &m.CreatedAt, &m.CreatedBy,
			&m.Player1.ID, &m.Player1.GoogleID, &m.Player1.Email, &m.Player1.Name, &m.Player1.AvatarURL, &m.Player1.MatchesPlayed, &m.Player1.MatchesWon, &m.Player1.CreatedAt, &m.Player1.UpdatedAt,
			&m.Player2.ID, &m.Player2.GoogleID, &m.Player2.Email, &m.Player2.Name, &m.Player2.AvatarURL, &m.Player2.MatchesPlayed, &m.Player2.MatchesWon, &m.Player2.CreatedAt, &m.Player2.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		matches = append(matches, m)
	}

	// Fetch sets for each match
	for i := range matches {
		sets, err := r.GetMatchSets(ctx, matches[i].ID)
		if err != nil {
			return nil, 0, err
		}
		matches[i].Sets = sets
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

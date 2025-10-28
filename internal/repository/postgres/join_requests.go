package postgres

import (
	"context"
	"database/sql"
	"go-data-catalog/internal/models"
)

type JoinRequestRepository struct {
	db *DB
}

func NewJoinRequestRepository(db *DB) *JoinRequestRepository { return &JoinRequestRepository{db: db} }

func (r *JoinRequestRepository) Create(ctx context.Context, teamID, userID int) (*models.JoinRequest, error) {
	jr := &models.JoinRequest{TeamID: teamID, UserID: userID, Status: "pending"}
	query := `
		INSERT INTO join_requests (team_id, user_id, status)
		VALUES ($1, $2, 'pending')
		RETURNING id, status, created_at
	`
	if err := r.db.Pool.QueryRow(ctx, query, teamID, userID).Scan(&jr.ID, &jr.Status, &jr.CreatedAt); err != nil {
		return nil, err
	}
	return jr, nil
}

func (r *JoinRequestRepository) GetByID(ctx context.Context, id int) (*models.JoinRequest, error) {
	query := `SELECT id, team_id, user_id, status, created_at, processed_by, processed_at FROM join_requests WHERE id = $1`
	var jr models.JoinRequest
	var pb sql.NullInt32
	var pa sql.NullTime
	if err := r.db.Pool.QueryRow(ctx, query, id).Scan(&jr.ID, &jr.TeamID, &jr.UserID, &jr.Status, &jr.CreatedAt, &pb, &pa); err != nil {
		return nil, err
	}
	if pb.Valid { v := int(pb.Int32); jr.ProcessedBy = &v }
	if pa.Valid { t := pa.Time; jr.ProcessedAt = &t }
	return &jr, nil
}

func (r *JoinRequestRepository) ListByTeam(ctx context.Context, teamID int, status string) ([]models.JoinRequest, error) {
	q := `SELECT id, team_id, user_id, status, created_at, processed_by, processed_at FROM join_requests WHERE team_id = $1`
	args := []any{teamID}
	if status != "" {
		q += " AND status = $2"
		args = append(args, status)
	}
	q += " ORDER BY created_at DESC"
	rows, err := r.db.Pool.Query(ctx, q, args...)
	if err != nil { return nil, err }
	defer rows.Close()
	var res []models.JoinRequest
	for rows.Next() {
		var jr models.JoinRequest
		var pb sql.NullInt32
		var pa sql.NullTime
		if err := rows.Scan(&jr.ID, &jr.TeamID, &jr.UserID, &jr.Status, &jr.CreatedAt, &pb, &pa); err != nil { return nil, err }
		if pb.Valid { v := int(pb.Int32); jr.ProcessedBy = &v }
		if pa.Valid { t := pa.Time; jr.ProcessedAt = &t }
		res = append(res, jr)
	}
	return res, nil
}

func (r *JoinRequestRepository) UpdateStatus(ctx context.Context, id int, processedBy int, status string) error {
	query := `
		UPDATE join_requests
		SET status = $2, processed_by = $3, processed_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.Pool.Exec(ctx, query, id, status, processedBy)
	return err
}

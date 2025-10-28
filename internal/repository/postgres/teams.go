package postgres

import (
	"context"
	"go-data-catalog/internal/models"
)

type TeamRepository struct {
	db *DB
}

func NewTeamRepository(db *DB) *TeamRepository { return &TeamRepository{db: db} }

func (r *TeamRepository) CreateTeam(ctx context.Context, t *models.Team) error {
	query := `
		INSERT INTO teams (name, description, created_by)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	return r.db.Pool.QueryRow(ctx, query, t.Name, t.Description, t.CreatedBy).Scan(&t.ID, &t.CreatedAt)
}

func (r *TeamRepository) GetByID(ctx context.Context, id int) (*models.Team, error) {
	query := `SELECT id, name, description, created_by, created_at FROM teams WHERE id = $1`
	var t models.Team
	if err := r.db.Pool.QueryRow(ctx, query, id).Scan(&t.ID, &t.Name, &t.Description, &t.CreatedBy, &t.CreatedAt); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TeamRepository) Search(ctx context.Context, q string, limit int) ([]models.Team, error) {
	if limit <= 0 || limit > 50 { limit = 20 }
	query := `
		SELECT id, name, description, created_by, created_at
		FROM teams
		WHERE name ILIKE '%' || $1 || '%'
		ORDER BY name ASC
		LIMIT $2
	`
	rows, err := r.db.Pool.Query(ctx, query, q, limit)
	if err != nil { return nil, err }
	defer rows.Close()
	var res []models.Team
	for rows.Next() {
		var t models.Team
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.CreatedBy, &t.CreatedAt); err != nil { return nil, err }
		res = append(res, t)
	}
	return res, nil
}

func (r *TeamRepository) ListForUser(ctx context.Context, userID int) ([]models.Team, error) {
	query := `
		SELECT t.id, t.name, t.description, t.created_by, t.created_at
		FROM teams t
		JOIN team_members m ON m.team_id = t.id AND m.user_id = $1 AND m.status = 'active'
		ORDER BY t.name
	`
	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil { return nil, err }
	defer rows.Close()
	var res []models.Team
	for rows.Next() {
		var t models.Team
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.CreatedBy, &t.CreatedAt); err != nil { return nil, err }
		res = append(res, t)
	}
	return res, nil
}

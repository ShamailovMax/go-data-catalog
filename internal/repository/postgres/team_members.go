package postgres

import (
	"context"
)

type TeamMemberRepository struct {
	db *DB
}

func NewTeamMemberRepository(db *DB) *TeamMemberRepository { return &TeamMemberRepository{db: db} }

func (r *TeamMemberRepository) AddOrUpdate(ctx context.Context, teamID, userID int, role string) error {
	query := `
		INSERT INTO team_members (team_id, user_id, role, status)
		VALUES ($1, $2, $3, 'active')
		ON CONFLICT (team_id, user_id)
		DO UPDATE SET role = EXCLUDED.role, status = 'active'
	`
	_, err := r.db.Pool.Exec(ctx, query, teamID, userID, role)
	return err
}

func (r *TeamMemberRepository) GetUserRole(ctx context.Context, teamID, userID int) (string, error) {
	query := `SELECT role FROM team_members WHERE team_id = $1 AND user_id = $2 AND status = 'active'`
	var role string
	if err := r.db.Pool.QueryRow(ctx, query, teamID, userID).Scan(&role); err != nil {
		return "", err
	}
	return role, nil
}

func (r *TeamMemberRepository) IsMember(ctx context.Context, teamID, userID int) (bool, error) {
	var exists bool
	if err := r.db.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM team_members WHERE team_id=$1 AND user_id=$2 AND status='active')`, teamID, userID).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

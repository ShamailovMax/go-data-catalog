package postgres

import (
	"context"
	"go-data-catalog/internal/models"
)

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, u *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, name, system_role, is_active)
		VALUES ($1, $2, $3, COALESCE($4, 'user'), COALESCE($5, TRUE))
		RETURNING id, system_role, is_active, created_at
	`
	return r.db.Pool.QueryRow(ctx, query, u.Email, u.PasswordHash, u.Name, u.SystemRole, u.IsActive).
		Scan(&u.ID, &u.SystemRole, &u.IsActive, &u.CreatedAt)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, email, password_hash, name, system_role, is_active, created_at FROM users WHERE email = $1`
	var u models.User
	if err := r.db.Pool.QueryRow(ctx, query, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.SystemRole, &u.IsActive, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `SELECT id, email, password_hash, name, system_role, is_active, created_at FROM users WHERE id = $1`
	var u models.User
	if err := r.db.Pool.QueryRow(ctx, query, id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.SystemRole, &u.IsActive, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

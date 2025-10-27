package postgres

import (
	"context"
	"go-data-catalog/internal/models"
)

type ArtifactFieldRepository struct {
	db *DB
}

func NewArtifactFieldRepository(db *DB) *ArtifactFieldRepository {
	return &ArtifactFieldRepository{db: db}
}

func (r *ArtifactFieldRepository) GetFieldsByArtifactID(ctx context.Context, artifactID int) ([]models.ArtifactField, error) {
	query := `
		SELECT id, artifact_id, field_name, data_type, description, is_pk, created_at
		FROM artifact_fields
		WHERE artifact_id = $1
		ORDER BY id
	`
	rows, err := r.db.Pool.Query(ctx, query, artifactID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fields []models.ArtifactField
	for rows.Next() {
		var f models.ArtifactField
		if err := rows.Scan(
			&f.ID,
			&f.ArtifactID,
			&f.FieldName,
			&f.DataType,
			&f.Description,
			&f.IsPK,
			&f.CreatedAt,
		); err != nil {
			return nil, err
		}
		fields = append(fields, f)
	}
	return fields, nil
}

func (r *ArtifactFieldRepository) GetFieldByID(ctx context.Context, id int) (*models.ArtifactField, error) {
	query := `
		SELECT id, artifact_id, field_name, data_type, description, is_pk, created_at
		FROM artifact_fields
		WHERE id = $1
	`
	var f models.ArtifactField
	if err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&f.ID,
		&f.ArtifactID,
		&f.FieldName,
		&f.DataType,
		&f.Description,
		&f.IsPK,
		&f.CreatedAt,
	); err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *ArtifactFieldRepository) CreateField(ctx context.Context, f *models.ArtifactField) error {
	query := `
		INSERT INTO artifact_fields (artifact_id, field_name, data_type, description, is_pk)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	return r.db.Pool.QueryRow(
		ctx,
		query,
		f.ArtifactID,
		f.FieldName,
		f.DataType,
		f.Description,
		f.IsPK,
	).Scan(&f.ID, &f.CreatedAt)
}

func (r *ArtifactFieldRepository) UpdateField(ctx context.Context, id int, f *models.ArtifactField) error {
	query := `
		UPDATE artifact_fields
		SET field_name = $2, data_type = $3, description = $4, is_pk = $5
		WHERE id = $1
		RETURNING artifact_id, created_at
	`
	if err := r.db.Pool.QueryRow(
		ctx,
		query,
		id,
		f.FieldName,
		f.DataType,
		f.Description,
		f.IsPK,
	).Scan(&f.ArtifactID, &f.CreatedAt); err != nil {
		return err
	}
	f.ID = id
	return nil
}

func (r *ArtifactFieldRepository) DeleteField(ctx context.Context, id int) error {
	_, err := r.db.Pool.Exec(ctx, `DELETE FROM artifact_fields WHERE id = $1`, id)
	return err
}

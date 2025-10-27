package postgres

import (
	"context"
	"go-data-catalog/internal/models"
)

type ArtifactRepository struct {
	db *DB
}

func NewArtifactRepository(db *DB) *ArtifactRepository{
	return &ArtifactRepository{db: db}
}

func (r *ArtifactRepository) GetAllArtifacts(ctx context.Context) ([]models.Artifact, error){
	query := `
		select id, name, type, description, project_name, developer_id, created_at
		from artifacts
		order by created_at desc
	`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var artifacts []models.Artifact
	for rows.Next() {
		var artifact models.Artifact
		err := rows.Scan(
			&artifact.ID,
			&artifact.Name, 
			&artifact.Type,
			&artifact.Description,
			&artifact.ProjectName,
			&artifact.DeveloperID,
			&artifact.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		artifacts = append(artifacts, artifact)
	}
	return artifacts, nil
}

func (r *ArtifactRepository) CreateArtifact(ctx context.Context, artifact *models.Artifact) error {
	query := `
		INSERT INTO artifacts (name, type, description, project_name, developer_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	err := r.db.Pool.QueryRow(
		ctx,
		query,
		artifact.Name,
		artifact.Type,
		artifact.Description, 
		artifact.ProjectName,
		artifact.DeveloperID,
	).Scan(&artifact.ID, &artifact.CreatedAt)
	
	return err
}

func (r *ArtifactRepository) GetArtifactByID(ctx context.Context, id int) (*models.Artifact, error) {
	query := `
		SELECT id, name, type, description, project_name, developer_id, created_at
		FROM artifacts
		WHERE id = $1
	`
	
	var artifact models.Artifact
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&artifact.ID,
		&artifact.Name,
		&artifact.Type,
		&artifact.Description,
		&artifact.ProjectName,
		&artifact.DeveloperID,
		&artifact.CreatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &artifact, nil
}

func (r *ArtifactRepository) UpdateArtifact(ctx context.Context, id int, artifact *models.Artifact) error {
	query := `
		UPDATE artifacts 
		SET name = $2, type = $3, description = $4, project_name = $5, developer_id = $6
		WHERE id = $1
		RETURNING created_at
	`
	
	err := r.db.Pool.QueryRow(
		ctx,
		query,
		id,
		artifact.Name,
		artifact.Type,
		artifact.Description,
		artifact.ProjectName,
		artifact.DeveloperID,
	).Scan(&artifact.CreatedAt)
	
	if err != nil {
		return err
	}
	
	artifact.ID = id
	return nil
}

func (r *ArtifactRepository) DeleteArtifact(ctx context.Context, id int) error {
	query := `DELETE FROM artifacts WHERE id = $1`
	
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}

// Exists проверяет, существует ли артефакт
func (r *ArtifactRepository) Exists(ctx context.Context, id int) (bool, error) {
	var exists bool
	err := r.db.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM artifacts WHERE id = $1)`, id).Scan(&exists)
	return exists, err
}

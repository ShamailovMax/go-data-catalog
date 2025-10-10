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

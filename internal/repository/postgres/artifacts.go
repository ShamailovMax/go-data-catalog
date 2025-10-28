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

func (r *ArtifactRepository) GetAllArtifacts(ctx context.Context, teamID int) ([]models.Artifact, error){
	query := `
		SELECT id, name, type, description, project_name, developer_id, team_id, created_at
		FROM artifacts
		WHERE team_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, teamID)
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
			&artifact.TeamID,
			&artifact.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		artifacts = append(artifacts, artifact)
	}
	return artifacts, nil
}

func (r *ArtifactRepository) CreateArtifact(ctx context.Context, teamID int, artifact *models.Artifact) error {
	query := `
		INSERT INTO artifacts (name, type, description, project_name, developer_id, team_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, team_id, created_at
	`

	err := r.db.Pool.QueryRow(
		ctx,
		query,
		artifact.Name,
		artifact.Type,
		artifact.Description, 
		artifact.ProjectName,
		artifact.DeveloperID,
		teamID,
	).Scan(&artifact.ID, &artifact.TeamID, &artifact.CreatedAt)
	
	return err
}

func (r *ArtifactRepository) GetArtifactByID(ctx context.Context, teamID, id int) (*models.Artifact, error) {
	query := `
		SELECT id, name, type, description, project_name, developer_id, team_id, created_at
		FROM artifacts
		WHERE id = $1 AND team_id = $2
	`
	
	var artifact models.Artifact
	err := r.db.Pool.QueryRow(ctx, query, id, teamID).Scan(
		&artifact.ID,
		&artifact.Name,
		&artifact.Type,
		&artifact.Description,
		&artifact.ProjectName,
		&artifact.DeveloperID,
		&artifact.TeamID,
		&artifact.CreatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &artifact, nil
}

func (r *ArtifactRepository) UpdateArtifact(ctx context.Context, teamID, id int, artifact *models.Artifact) error {
	query := `
		UPDATE artifacts 
		SET name = $3, type = $4, description = $5, project_name = $6, developer_id = $7
		WHERE id = $1 AND team_id = $2
		RETURNING team_id, created_at
	`
	
	err := r.db.Pool.QueryRow(
		ctx,
		query,
		id,
		teamID,
		artifact.Name,
		artifact.Type,
		artifact.Description,
		artifact.ProjectName,
		artifact.DeveloperID,
	).Scan(&artifact.TeamID, &artifact.CreatedAt)
	
	if err != nil {
		return err
	}
	
	artifact.ID = id
	return nil
}

func (r *ArtifactRepository) DeleteArtifact(ctx context.Context, teamID, id int) error {
	query := `DELETE FROM artifacts WHERE id = $1 AND team_id = $2`
	
	_, err := r.db.Pool.Exec(ctx, query, id, teamID)
	return err
}

// Exists проверяет, существует ли артефакт
func (r *ArtifactRepository) Exists(ctx context.Context, teamID, id int) (bool, error) {
	var exists bool
	err := r.db.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM artifacts WHERE id = $1 AND team_id = $2)`, id, teamID).Scan(&exists)
	return exists, err
}

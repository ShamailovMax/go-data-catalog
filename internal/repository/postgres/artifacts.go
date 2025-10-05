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
	return []models.Artifact{}, nil
}

func (r *ArtifactRepository) CreateArtifact(ctx context.Context, artifact *models.Artifact) error {
	return nil
}

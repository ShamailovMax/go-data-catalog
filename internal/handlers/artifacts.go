package handlers

import (
	"net/http"
	"go-data-catalog/internal/models"
	"go-data-catalog/internal/repository/postgres"

	"github.com/gin-gonic/gin"
)

type ArtifactHandler struct {
	repo *postgres.ArtifactRepository
}

func NewArtifactHandler(repo *postgres.ArtifactRepository) *ArtifactHandler {
	return &ArtifactHandler{repo: repo}
}

func (h *ArtifactHandler) GetArtifacts(c *gin.Context) {
	artifacts, err := h.repo.GetAllArtifacts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, artifacts)
}

func (h *ArtifactHandler) CreateArtifact(c *gin.Context) {
	var artifact models.Artifact
	if err := c.BindJSON(&artifact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	
	if err := h.repo.CreateArtifact(c.Request.Context(), &artifact); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, artifact)
}
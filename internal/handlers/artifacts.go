package handlers

import (
	"net/http"
	"strconv"
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

func (h *ArtifactHandler) teamID(c *gin.Context) (int, bool) {
	teamIDParam := c.Param("teamId")
	teamID, err := strconv.Atoi(teamIDParam)
	if err != nil || teamID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return 0, false
	}
	return teamID, true
}

func (h *ArtifactHandler) GetArtifacts(c *gin.Context) {
	teamID, ok := h.teamID(c); if !ok { return }
	artifacts, err := h.repo.GetAllArtifacts(c.Request.Context(), teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, artifacts)
}

func (h *ArtifactHandler) CreateArtifact(c *gin.Context) {
	teamID, ok := h.teamID(c); if !ok { return }
	var artifact models.Artifact
	if err := c.BindJSON(&artifact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	
	if err := h.repo.CreateArtifact(c.Request.Context(), teamID, &artifact); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, artifact)
}

func (h *ArtifactHandler) GetArtifactByID(c *gin.Context) {
	teamID, ok := h.teamID(c); if !ok { return }
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	artifact, err := h.repo.GetArtifactByID(c.Request.Context(), teamID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artifact not found"})
		return
	}
	
	c.JSON(http.StatusOK, artifact)
}

func (h *ArtifactHandler) UpdateArtifact(c *gin.Context) {
	teamID, ok := h.teamID(c); if !ok { return }
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	var artifact models.Artifact
	if err := c.BindJSON(&artifact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	
	if err := h.repo.UpdateArtifact(c.Request.Context(), teamID, id, &artifact); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, artifact)
}

func (h *ArtifactHandler) DeleteArtifact(c *gin.Context) {
	teamID, ok := h.teamID(c); if !ok { return }
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := h.repo.DeleteArtifact(c.Request.Context(), teamID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Artifact deleted successfully"})
}

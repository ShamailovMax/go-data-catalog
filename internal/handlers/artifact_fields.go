package handlers

import (
	"net/http"
	"strconv"
	"go-data-catalog/internal/models"
	"go-data-catalog/internal/repository/postgres"

	"github.com/gin-gonic/gin"
)

type ArtifactFieldHandler struct {
	repo         *postgres.ArtifactFieldRepository
	artifactRepo *postgres.ArtifactRepository
}

func NewArtifactFieldHandler(repo *postgres.ArtifactFieldRepository, artifactRepo *postgres.ArtifactRepository) *ArtifactFieldHandler {
	return &ArtifactFieldHandler{repo: repo, artifactRepo: artifactRepo}
}

func (h *ArtifactFieldHandler) teamID(c *gin.Context) (int, bool) {
	teamIDParam := c.Param("teamId")
	teamID, err := strconv.Atoi(teamIDParam)
	if err != nil || teamID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return 0, false
	}
	return teamID, true
}

// List fields for specific artifact
func (h *ArtifactFieldHandler) GetFieldsByArtifact(c *gin.Context) {
	teamID, ok := h.teamID(c); if !ok { return }
	artifactIDParam := c.Param("id")
	artifactID, err := strconv.Atoi(artifactIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid artifact_id"})
		return
	}

	// ensure artifact belongs to the team
	if ok, _ := h.artifactExists(c, teamID, artifactID); !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artifact not found"})
		return
	}

	fields, err := h.repo.GetFieldsByArtifactID(c.Request.Context(), artifactID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fields)
}

// Create field under artifact
func (h *ArtifactFieldHandler) CreateField(c *gin.Context) {
	teamID, ok := h.teamID(c); if !ok { return }
	artifactIDParam := c.Param("id")
	artifactID, err := strconv.Atoi(artifactIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid artifact_id"})
		return
	}

	// Проверяем, что артефакт существует в рамках команды
	exists, err := h.artifactRepo.Exists(c.Request.Context(), teamID, artifactID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artifact not found"})
		return
	}

	var f models.ArtifactField
	if err := c.BindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	// Принудительно используем artifact_id из пути
	f.ArtifactID = artifactID

	if err := h.repo.CreateField(c.Request.Context(), &f); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, f)
}

func (h *ArtifactFieldHandler) GetFieldByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	f, err := h.repo.GetFieldByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Field not found"})
		return
	}
	c.JSON(http.StatusOK, f)
}

func (h *ArtifactFieldHandler) UpdateField(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var f models.ArtifactField
	if err := c.BindJSON(&f); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.repo.UpdateField(c.Request.Context(), id, &f); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, f)
}

func (h *ArtifactFieldHandler) DeleteField(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.repo.DeleteField(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Field deleted successfully"})
}

func (h *ArtifactFieldHandler) artifactExists(c *gin.Context, teamID, artifactID int) (bool, error) {
	return h.artifactRepo.Exists(c.Request.Context(), teamID, artifactID)
}

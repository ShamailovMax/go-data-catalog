package handlers

import (
	"net/http"
	"strconv"
	"go-data-catalog/internal/models"
	"go-data-catalog/internal/repository/postgres"

	"github.com/gin-gonic/gin"
)

type ContactHandler struct {
	repo *postgres.ContactRepository
}

func NewContactHandler(repo *postgres.ContactRepository) *ContactHandler {
	return &ContactHandler{repo: repo}
}

func (h *ContactHandler) teamID(c *gin.Context) (int, bool) {
	teamIDParam := c.Param("teamId")
	teamID, err := strconv.Atoi(teamIDParam)
	if err != nil || teamID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
		return 0, false
	}
	return teamID, true
}

func (h *ContactHandler) GetContacts(c *gin.Context) {
	teamID, ok := h.teamID(c); if !ok { return }
	contacts, err := h.repo.GetAllContacts(c.Request.Context(), teamID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, contacts)
}

func (h *ContactHandler) GetContactByID(c *gin.Context) {
	teamID, ok := h.teamID(c); if !ok { return }
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	contact, err := h.repo.GetContactByID(c.Request.Context(), teamID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}
	
	c.JSON(http.StatusOK, contact)
}

func (h *ContactHandler) CreateContact(c *gin.Context) {
	teamID, ok := h.teamID(c); if !ok { return }
	var contact models.Contact
	if err := c.BindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	
	if err := h.repo.CreateContact(c.Request.Context(), teamID, &contact); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, contact)
}

func (h *ContactHandler) UpdateContact(c *gin.Context) {
	teamID, ok := h.teamID(c); if !ok { return }
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	var contact models.Contact
	if err := c.BindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	
	if err := h.repo.UpdateContact(c.Request.Context(), teamID, id, &contact); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, contact)
}

func (h *ContactHandler) DeleteContact(c *gin.Context) {
	teamID, ok := h.teamID(c); if !ok { return }
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := h.repo.DeleteContact(c.Request.Context(), teamID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Contact deleted successfully"})
}

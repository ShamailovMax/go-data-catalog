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

func (h *ContactHandler) GetContacts(c *gin.Context) {
	contacts, err := h.repo.GetAllContacts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, contacts)
}

func (h *ContactHandler) GetContactByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	contact, err := h.repo.GetContactByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}
	
	c.JSON(http.StatusOK, contact)
}

func (h *ContactHandler) CreateContact(c *gin.Context) {
	var contact models.Contact
	if err := c.BindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	
	if err := h.repo.CreateContact(c.Request.Context(), &contact); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, contact)
}

func (h *ContactHandler) UpdateContact(c *gin.Context) {
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
	
	if err := h.repo.UpdateContact(c.Request.Context(), id, &contact); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, contact)
}

func (h *ContactHandler) DeleteContact(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := h.repo.DeleteContact(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Contact deleted successfully"})
}
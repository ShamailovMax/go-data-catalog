package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"go-data-catalog/internal/config"
	"go-data-catalog/internal/models"
	"go-data-catalog/internal/repository/postgres"
)

type AuthHandler struct {
	users *postgres.UserRepository
	cfg   *config.Config
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	Name     string `json:"name" binding:"omitempty,max=255"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type authResponse struct {
	Token string       `json:"token"`
	User  models.User  `json:"user"`
}

func NewAuthHandler(users *postgres.UserRepository, cfg *config.Config) *AuthHandler {
	return &AuthHandler{users: users, cfg: cfg}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}
	u := &models.User{Email: req.Email, PasswordHash: string(hash), Name: req.Name, SystemRole: "user", IsActive: true}
	if err := h.users.CreateUser(c.Request.Context(), u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists or invalid"})
		return
	}
	token, err := h.issueToken(u)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to issue token"}); return }
	u.PasswordHash = ""
	c.JSON(http.StatusCreated, authResponse{Token: token, User: *u})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	u, err := h.users.GetByEmail(c.Request.Context(), req.Email)
	if err != nil { c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"}); return }
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil { c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"}); return }
	token, err := h.issueToken(u)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to issue token"}); return }
	u.PasswordHash = ""
	c.JSON(http.StatusOK, authResponse{Token: token, User: *u})
}

func (h *AuthHandler) issueToken(u *models.User) (string, error) {
	expires := time.Now().Add(time.Duration(h.cfg.TokenTTL) * time.Minute)
	claims := jwt.MapClaims{
		"user_id": u.ID,
		"system_role": u.SystemRole,
		"exp": expires.Unix(),
		"iat": time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(h.cfg.JWTSecret))
}

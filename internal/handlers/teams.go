package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go-data-catalog/internal/middleware"
	"go-data-catalog/internal/models"
	"go-data-catalog/internal/repository/postgres"
)

type TeamsHandler struct {
	teams      *postgres.TeamRepository
	members    *postgres.TeamMemberRepository
	joinReqs   *postgres.JoinRequestRepository
}

type createTeamRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=255"`
	Description string `json:"description" binding:"omitempty,max=1000"`
}

type joinDecisionRequest struct {
	Action string `json:"action" binding:"required,oneof=approve reject"`
}

func NewTeamsHandler(teams *postgres.TeamRepository, members *postgres.TeamMemberRepository, joinReqs *postgres.JoinRequestRepository) *TeamsHandler {
	return &TeamsHandler{teams: teams, members: members, joinReqs: joinReqs}
}

// POST /api/v1/teams
func (h *TeamsHandler) CreateTeam(c *gin.Context) {
	userID := c.GetInt(middleware.CtxUserID)
	var req createTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"}); return }
	t := &models.Team{Name: req.Name, Description: req.Description, CreatedBy: userID}
	if err := h.teams.CreateTeam(c.Request.Context(), t); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "cannot create team"}); return }
	_ = h.members.AddOrUpdate(c.Request.Context(), t.ID, userID, "owner")
	c.JSON(http.StatusCreated, t)
}

// GET /api/v1/teams?search=foo
func (h *TeamsHandler) Search(c *gin.Context) {
	q := c.Query("search")
	res, err := h.teams.Search(c.Request.Context(), q, 20)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"}); return }
	c.JSON(http.StatusOK, res)
}

// POST /api/v1/teams/:teamId/join
func (h *TeamsHandler) RequestJoin(c *gin.Context) {
	userID := c.GetInt(middleware.CtxUserID)
	teamID, _ := strconv.Atoi(c.Param("teamId"))
	isMember, err := h.members.IsMember(c.Request.Context(), teamID, userID)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"}); return }
	if isMember { c.JSON(http.StatusBadRequest, gin.H{"error": "already a member"}); return }
	jr, err := h.joinReqs.Create(c.Request.Context(), teamID, userID)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "cannot create request"}); return }
	c.JSON(http.StatusCreated, jr)
}

// GET /api/v1/teams/:teamId/requests?status=pending (owner/admin)
func (h *TeamsHandler) ListRequests(c *gin.Context) {
	teamID, _ := strconv.Atoi(c.Param("teamId"))
	status := c.DefaultQuery("status", "")
	items, err := h.joinReqs.ListByTeam(c.Request.Context(), teamID, status)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"}); return }
	c.JSON(http.StatusOK, items)
}

// POST /api/v1/teams/:teamId/requests/:id/approve or reject
func (h *TeamsHandler) DecideRequest(c *gin.Context) {
	teamID, _ := strconv.Atoi(c.Param("teamId"))
	reqID, _ := strconv.Atoi(c.Param("id"))
	actorID := c.GetInt(middleware.CtxUserID)
	action := c.Param("action") // expects "approve" or "reject"
	if action != "approve" && action != "reject" { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid action"}); return }
	jr, err := h.joinReqs.GetByID(c.Request.Context(), reqID)
	if err != nil || jr.TeamID != teamID { c.JSON(http.StatusNotFound, gin.H{"error": "request not found"}); return }
	status := map[string]string{"approve": "approved", "reject": "rejected"}[action]
	if err := h.joinReqs.UpdateStatus(c.Request.Context(), reqID, actorID, status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot update"}); return
	}
	if status == "approved" {
		_ = h.members.AddOrUpdate(c.Request.Context(), teamID, jr.UserID, "member")
	}
	c.Status(http.StatusNoContent)
}

// GET /api/v1/me/teams
func (h *TeamsHandler) MyTeams(c *gin.Context) {
	userID := c.GetInt(middleware.CtxUserID)
	items, err := h.teams.ListForUser(c.Request.Context(), userID)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"}); return }
	c.JSON(http.StatusOK, items)
}

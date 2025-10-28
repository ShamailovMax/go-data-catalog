package middleware

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"go-data-catalog/internal/repository/postgres"
)

type ctxKey string

const (
	CtxUserID   = "userID"
	CtxTeamID   = "teamID"
	CtxTeamRole = "teamRole"
	CtxSysRole  = "systemRole"
)

// TeamMembershipMiddleware ensures the authenticated user is a member of the team in path param :teamId
func TeamMembershipMiddleware(membersRepo *postgres.TeamMemberRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDVal, exists := c.Get(CtxUserID)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		userID := userIDVal.(int)

		teamIDStr := c.Param("teamId")
		teamID, err := strconv.Atoi(teamIDStr)
		if err != nil || teamID <= 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid team id"})
			return
		}
		isMember, err := membersRepo.IsMember(c.Request.Context(), teamID, userID)
		if err != nil || !isMember {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "not a team member"})
			return
		}
		role, _ := membersRepo.GetUserRole(c.Request.Context(), teamID, userID)
		c.Set(CtxTeamID, teamID)
		c.Set(CtxTeamRole, role)
		c.Next()
	}
}

// RequireTeamRole allows only specified roles
func RequireTeamRole(roles ...string) gin.HandlerFunc {
	allowed := map[string]struct{}{}
	for _, r := range roles { allowed[r] = struct{}{} }
	return func(c *gin.Context) {
		val, ok := c.Get(CtxTeamRole)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		role := val.(string)
		if _, ok := allowed[role]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}

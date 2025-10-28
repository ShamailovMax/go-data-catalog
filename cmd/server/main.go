package main

import (
	"log"
	"go-data-catalog/internal/config"
	"go-data-catalog/internal/handlers"
	"go-data-catalog/internal/middleware"
	"go-data-catalog/internal/repository/postgres"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	
	// Инициализация БД
	db, err := postgres.NewDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	
	// Инициализация репозиториев
	artifactRepo := postgres.NewArtifactRepository(db)
	contactRepo := postgres.NewContactRepository(db)
	artifactFieldRepo := postgres.NewArtifactFieldRepository(db)
	userRepo := postgres.NewUserRepository(db)
	teamRepo := postgres.NewTeamRepository(db)
	memberRepo := postgres.NewTeamMemberRepository(db)
	joinReqRepo := postgres.NewJoinRequestRepository(db)
	
	// Инициализация handlers
	artifactHandler := handlers.NewArtifactHandler(artifactRepo)
	contactHandler := handlers.NewContactHandler(contactRepo)
	artifactFieldHandler := handlers.NewArtifactFieldHandler(artifactFieldRepo, artifactRepo)
	authHandler := handlers.NewAuthHandler(userRepo, cfg)
	teamsHandler := handlers.NewTeamsHandler(teamRepo, memberRepo, joinReqRepo)
	
	// Настройка роутера
	r := gin.New() // Используем New вместо Default чтобы сами настроить middleware
	
	// Добавляем наши middleware
	r.Use(gin.Recovery()) // Восстановление после паники
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.ErrorHandlerMiddleware())
	r.Use(middleware.CORSMiddleware())
	
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	// Static files (frontend)
	r.Static("/static", "./web/static")
	r.GET("/", func(c *gin.Context) {
		c.File("./web/static/index.html")
	})
	
	// Public auth
	v1 := r.Group("/api/v1")
	{
		v1.POST("/auth/register", authHandler.Register)
		v1.POST("/auth/login", authHandler.Login)
	}
	
	// Authenticated routes
	v1auth := r.Group("/api/v1")
	v1auth.Use(middleware.AuthMiddleware(cfg))
	{
		// teams discovery/creation
		v1auth.GET("/teams", teamsHandler.Search)
		v1auth.POST("/teams", teamsHandler.CreateTeam)
		v1auth.POST("/teams/:teamId/join", teamsHandler.RequestJoin)
		v1auth.GET("/me/teams", teamsHandler.MyTeams)

		// team-scoped routes
		team := v1auth.Group("/teams/:teamId")
		team.Use(middleware.TeamMembershipMiddleware(memberRepo))
		{
			// join request admin endpoints (owner/admin only)
			admin := team.Group("")
			admin.Use(middleware.RequireTeamRole("owner", "admin"))
			{
				admin.GET("/requests", teamsHandler.ListRequests)
				admin.POST("/requests/:id/:action", teamsHandler.DecideRequest) // action=approve|reject
			}

			// artifacts
			artifacts := team.Group("/artifacts")
			{
				artifacts.GET("", artifactHandler.GetArtifacts)
				artifacts.GET("/:id", artifactHandler.GetArtifactByID)
				artifacts.POST("", artifactHandler.CreateArtifact)
				artifacts.PUT("/:id", artifactHandler.UpdateArtifact)
				artifacts.DELETE("/:id", artifactHandler.DeleteArtifact)
				// artifact fields
				artifactFields := artifacts.Group("/:id/fields")
				{
					artifactFields.GET("", artifactFieldHandler.GetFieldsByArtifact)
					artifactFields.POST("", artifactFieldHandler.CreateField)
				}
			}

			// contacts
			contacts := team.Group("/contacts")
			{
				contacts.GET("", contactHandler.GetContacts)
				contacts.GET("/:id", contactHandler.GetContactByID)
				contacts.POST("", contactHandler.CreateContact)
				contacts.PUT("/:id", contactHandler.UpdateContact)
				contacts.DELETE("/:id", contactHandler.DeleteContact)
			}

			// fields by id
			fields := team.Group("/fields")
			{
				fields.GET("/:id", artifactFieldHandler.GetFieldByID)
				fields.PUT("/:id", artifactFieldHandler.UpdateField)
				fields.DELETE("/:id", artifactFieldHandler.DeleteField)
			}
		}
	}
	
	log.Println("Server starting on :" + cfg.ServerPort)
	r.Run(":" + cfg.ServerPort)
}

package models

import "time"

type Contact struct {
    ID              int       `json:"id"`
    Name            string    `json:"name" binding:"required,min=2,max=255"`
    TelegramContact string    `json:"telegram_contact" binding:"omitempty,min=3,max=100"`
    TeamID          int       `json:"team_id"`
    CreatedAt       time.Time `json:"created_at"`
}

type Artifact struct {
    ID          int       `json:"id"`
    Name        string    `json:"name" binding:"required,min=2,max=255"`
    Type        string    `json:"type" binding:"required,oneof=table view procedure function index dataset api file"`
    Description string    `json:"description" binding:"omitempty,max=1000"`
    ProjectName string    `json:"project_name" binding:"required,min=2,max=255"`
    DeveloperID int       `json:"developer_id" binding:"omitempty,min=1"`
    TeamID      int       `json:"team_id"`
    CreatedAt   time.Time `json:"created_at"`
}

type ArtifactField struct {
    ID          int       `json:"id"`
    ArtifactID  int       `json:"artifact_id"`
    FieldName   string    `json:"field_name" binding:"required,min=1,max=255"`
    DataType    string    `json:"data_type" binding:"required,min=1,max=100"`
    Description string    `json:"description" binding:"omitempty,max=1000"`
    IsPK        bool      `json:"is_pk"`
    CreatedAt   time.Time `json:"created_at"`
}

// New auth/teams models

type User struct {
    ID           int       `json:"id"`
    Email        string    `json:"email" binding:"required,email"`
    PasswordHash string    `json:"-"`
    Name         string    `json:"name"`
    SystemRole   string    `json:"system_role"`
    IsActive     bool      `json:"is_active"`
    CreatedAt    time.Time `json:"created_at"`
}

type Team struct {
    ID          int       `json:"id"`
    Name        string    `json:"name" binding:"required,min=2,max=255"`
    Description string    `json:"description"`
    CreatedBy   int       `json:"created_by"`
    CreatedAt   time.Time `json:"created_at"`
}

type TeamMember struct {
    TeamID   int       `json:"team_id"`
    UserID   int       `json:"user_id"`
    Role     string    `json:"role"`
    Status   string    `json:"status"`
    JoinedAt time.Time `json:"joined_at"`
}

type JoinRequest struct {
    ID          int        `json:"id"`
    TeamID      int        `json:"team_id"`
    UserID      int        `json:"user_id"`
    Status      string     `json:"status"`
    CreatedAt   time.Time  `json:"created_at"`
    ProcessedBy *int       `json:"processed_by"`
    ProcessedAt *time.Time `json:"processed_at"`
}

package models

import "time"

type Contact struct {
    ID              int       `json:"id"`
    Name            string    `json:"name" binding:"required,min=2,max=255"`
    TelegramContact string    `json:"telegram_contact" binding:"omitempty,min=3,max=100"`
    CreatedAt       time.Time `json:"created_at"`
}

type Artifact struct {
    ID          int       `json:"id"`
    Name        string    `json:"name" binding:"required,min=2,max=255"`
    Type        string    `json:"type" binding:"required,oneof=table view procedure function index dataset api file"`
    Description string    `json:"description" binding:"omitempty,max=1000"`
    ProjectName string    `json:"project_name" binding:"required,min=2,max=255"`
    DeveloperID int       `json:"developer_id" binding:"omitempty,min=1"`
    CreatedAt   time.Time `json:"created_at"`
}

type ArtifactField struct {
    ID          int       `json:"id"`
    ArtifactID  int       `json:"artifact_id" binding:"required,min=1"`
    FieldName   string    `json:"field_name" binding:"required,min=1,max=255"`
    DataType    string    `json:"data_type" binding:"required,min=1,max=100"`
    Description string    `json:"description" binding:"omitempty,max=1000"`
    IsPK        bool      `json:"is_pk"`
    CreatedAt   time.Time `json:"created_at"`
}

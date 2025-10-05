package models

import "time"

type Contact struct {
    ID              int       `json:"id"`
    Name            string    `json:"name"`
    TelegramContact string    `json:"telegram_contact"`
    CreatedAt       time.Time `json:"created_at"`
}

type Artifact struct {
    ID          int       `json:"id"`
    Name        string    `json:"name"`
    Type        string    `json:"type"`
    Description string    `json:"description"`
    ProjectName string    `json:"project_name"`
    DeveloperID int       `json:"developer_id"`
    CreatedAt   time.Time `json:"created_at"`
}

type ArtifactField struct {
    ID          int    `json:"id"`
    ArtifactID  int    `json:"artifact_id"`
    FieldName   string `json:"field_name"`
    DataType    string `json:"data_type"`
    Description string `json:"description"`
    IsPK        bool   `json:"is_pk"`
}
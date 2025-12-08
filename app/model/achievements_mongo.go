package model

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Achievement mencerminkan dokumen di MongoDB
type Achievement struct {
	ID              primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	StudentID       uuid.UUID              `bson:"studentId" json:"student_id"`
	AchievementType string                 `bson:"achievementType" json:"achievement_type"`
	Title           string                 `bson:"title" json:"title"`
	Description     string                 `bson:"description" json:"description"`

	// Field Details dinamis (disimpan sebagai Map agar fleksibel sesuai SRS)
	Details map[string]interface{} `bson:"details" json:"details"`

	// Attachments (Array of struct)
	Attachments []Attachment `bson:"attachments" json:"attachments"`

	Tags      []string  `bson:"tags" json:"tags"`
	Points    int       `bson:"points" json:"points"`
	CreatedAt time.Time `bson:"createdAt" json:"created_at"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updated_at"`
}

type Attachment struct {
	FileName   string    `bson:"fileName" json:"file_name"`
	FileURL    string    `bson:"fileUrl" json:"file_url"`
	FileType   string    `bson:"fileType" json:"file_type"`
	UploadedAt time.Time `bson:"uploadedAt" json:"uploaded_at"`
}
package model

import "github.com/google/uuid"

// Response untuk /api/v1/reports/statistics
type SystemStatisticsResponse struct {
	TotalStudents     int            `json:"total_students"`
	TotalLecturers    int            `json:"total_lecturers"`
	TotalAchievements int            `json:"total_achievements"`
	ByStatus          map[string]int `json:"achievements_by_status"` // Draft, Submitted, Verified
	ByType            map[string]int `json:"achievements_by_type"`   // Competition, Seminar, etc (Dari Mongo)
}

// Response untuk /api/v1/reports/student/:id
type StudentReportResponse struct {
	StudentID         string                  `json:"student_id"`
	FullName          string                  `json:"full_name"`
	ProgramStudy      string                  `json:"program_study"`
	TotalAchievements int                     `json:"total_achievements"`
	TotalPoints       int                     `json:"total_points"`       // Hitung poin dari Mongo
	Achievements      []SimpleAchievementView `json:"achievements_list"`
}

type SimpleAchievementView struct {
	ID              uuid.UUID `json:"id"`
	Title           string    `json:"title"`
	AchievementType string    `json:"achievement_type"`
	Status          string    `json:"status"`
	Points          int       `json:"points"`
}
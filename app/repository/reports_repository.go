package repository

import (
	"context"
	"database/sql"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// --- POSTGRESQL STATS ---

func CountTotalUsersByRole(db *sql.DB, roleName string) (int, error) {
	// Asumsi kita hitung dari tabel users join roles, atau tabel students/lecturers langsung
	// Cara cepat: Hitung tabel students/lecturers
	var count int
	var query string
	if roleName == "student" {
		query = "SELECT COUNT(*) FROM students"
	} else {
		query = "SELECT COUNT(*) FROM lecturers"
	}
	err := db.QueryRow(query).Scan(&count)
	return count, err
}

func CountAchievementsByStatus(db *sql.DB) (map[string]int, error) {
	rows, err := db.Query("SELECT status, COUNT(*) FROM achievement_references GROUP BY status")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			continue
		}
		result[status] = count
	}
	return result, nil
}

// --- MONGODB STATS (AGGREGATION) ---

// Hitung jumlah prestasi berdasarkan Type (Competition, Organization, dll)
func AggregateAchievementsByType(db *mongo.Database) (map[string]int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.Collection("achievements")

	// PERBAIKAN: Menggunakan bson.M agar tidak kena linter error "unkeyed fields"
	pipeline := bson.A{
		bson.M{"$group": bson.M{
			"_id":   "$achievementType",
			"count": bson.M{"$sum": 1},
		}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	result := make(map[string]int)
	var items []struct {
		ID    string `bson:"_id"`
		Count int    `bson:"count"`
	}

	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	for _, item := range items {
		result[item.ID] = item.Count
	}
	return result, nil
}

// Hitung total poin mahasiswa tertentu (Sum Points)
func SumStudentPoints(db *mongo.Database, studentUUIDStr string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.Collection("achievements")

	// PERBAIKAN: Menggunakan bson.M untuk match dan group
	pipeline := bson.A{
		bson.M{"$match": bson.M{"studentId": studentUUIDStr}},
		bson.M{"$group": bson.M{
			"_id":         nil,
			"totalPoints": bson.M{"$sum": "$points"},
		}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		TotalPoints int `bson:"totalPoints"`
	}
	if err = cursor.All(ctx, &results); err != nil {
		return 0, err
	}

	if len(results) > 0 {
		return results[0].TotalPoints, nil
	}
	return 0, nil
}
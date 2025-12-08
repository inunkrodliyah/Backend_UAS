package repository

import (
	"context"
	"project-uas/app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName = "achievements"

// InsertAchievementMongo: Simpan ke Mongo -> Balikin ID String
func InsertAchievementMongo(db *mongo.Database, data model.Achievement) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Generate ID baru & Timestamp
	data.ID = primitive.NewObjectID()
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()

	// Insert
	res, err := db.Collection(collectionName).InsertOne(ctx, data)
	if err != nil {
		return "", err
	}

	// Ambil ID yang baru dibuat, konversi ke Hex String
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

// GetAchievementMongoByID: Ambil detail berdasarkan ID string hex
func GetAchievementMongoByID(db *mongo.Database, hexID string) (*model.Achievement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(hexID)
	if err != nil {
		return nil, err
	}

	var result model.Achievement
	err = db.Collection(collectionName).FindOne(ctx, bson.M{"_id": objID}).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
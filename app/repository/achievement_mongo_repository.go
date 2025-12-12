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

// Insert
func InsertAchievementMongo(db *mongo.Database, data model.Achievement) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	data.ID = primitive.NewObjectID()
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	res, err := db.Collection(collectionName).InsertOne(ctx, data)
	if err != nil { return "", err }
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

// Get By ID
func GetAchievementMongoByID(db *mongo.Database, hexID string) (*model.Achievement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(hexID)
	if err != nil { return nil, err }
	var result model.Achievement
	err = db.Collection(collectionName).FindOne(ctx, bson.M{"_id": objID}).Decode(&result)
	if err != nil { return nil, err }
	return &result, nil
}

// Update Data (Untuk Draft)
func UpdateAchievementMongo(db *mongo.Database, hexID string, data map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(hexID)
	if err != nil { return err }
	
	data["updatedAt"] = time.Now()
	_, err = db.Collection(collectionName).UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": data})
	return err
}

// Delete Data
func DeleteAchievementMongo(db *mongo.Database, hexID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(hexID)
	if err != nil { return err }
	_, err = db.Collection(collectionName).DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

// Add Attachment
func AddAttachmentMongo(db *mongo.Database, hexID string, attachment model.Attachment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	objID, err := primitive.ObjectIDFromHex(hexID)
	if err != nil { return err }
	
	update := bson.M{"$push": bson.M{"attachments": attachment}, "$set": bson.M{"updatedAt": time.Now()}}
	_, err = db.Collection(collectionName).UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}
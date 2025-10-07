package models

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
}

func InsertUser(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := mongoClient.Database(db).Collection(dbColumn)
	res, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Error inserting user: %v", err)
		return err
	}
	fmt.Println("Inserted user with ID: ", res.InsertedID)
	return err
}

func InsertUsers(users []User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := mongoClient.Database(db).Collection(dbColumn)
	newUsers := make([]interface{}, len(users))
	for i, user := range users {
		newUsers[i] = user
	}
	res, err := collection.InsertMany(ctx, newUsers)
	if err != nil {
		log.Printf("Error inserting users: %v", err)
		return err
	}
	fmt.Println("Inserted users with IDs: ", res.InsertedIDs)
	return nil
}

func UpdateUser(userId string, user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		return err
	}
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"username": user.Username, "email": user.Email, "password": user.Password}}

	collection := mongoClient.Database(db).Collection(dbColumn)
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return err
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", result.MatchedCount, result.ModifiedCount)
	return nil
}

func DeleteUser(userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		return err
	}
	filter := bson.M{"_id": id}
	collection := mongoClient.Database(db).Collection(dbColumn)
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		return err
	}
	fmt.Printf("Deleted %v documents in the users collection\n", result.DeletedCount)
	return nil
}

func FindUserById(userId string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result User
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		return User{}, fmt.Errorf("invalid user ID: %v", err)
	}
	filter := bson.M{"_id": id}
	collection := mongoClient.Database(db).Collection(dbColumn)
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		log.Printf("User not found for ID: %s", userId)
		return User{}, fmt.Errorf("user not found")
	}
	if err != nil {
		log.Printf("Error finding user: %v", err)
		return User{}, err
	}
	return result, nil
}

// Удаляем FindUser, так как она дублирует FindUserById
// Если FindAllUsers должна возвращать всех пользователей:
func ListAll() ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var results []User
	collection := mongoClient.Database(db).Collection(dbColumn)
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error occurred during Find(): %v", err)
		return nil, err
	}
	defer cur.Close(ctx)

	if err := cur.All(ctx, &results); err != nil {
		log.Printf("Error occurred during cursor.All(): %v", err)
		return nil, err
	}

	if len(results) == 0 {
		log.Println("No users found")
	}

	return results, nil
}

func DeleteAll() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := mongoClient.Database(db).Collection(dbColumn)
	result, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Printf("Error deleting all users: %v", err)
		return err
	}
	fmt.Printf("Deleted %v documents in the users collection\n", result.DeletedCount)
	return nil
}

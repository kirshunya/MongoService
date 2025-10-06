package models

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"log"
)

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
}

func InsertUser(user User) error {
	collection := mongoClient.Database(db).Collection(dbColumn)
	res, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted user with ID: ", res.InsertedID)
	return err
}

func InsertUsers(users []User) error {
	collection := mongoClient.Database(db).Collection(dbColumn)
	newUsers := make([]interface{}, len(users))
	for i, user := range users {
		newUsers[i] = user
	}
	res, err := collection.InsertMany(context.TODO(), newUsers)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted users with ID: ", res.InsertedIDs)
	return err
}

func UpdateUser(userId string, user User) error {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"username": user.Username, "email": user.Email, "password": user.Password}}

	collection := mongoClient.Database(db).Collection(dbColumn)
	result, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", result.MatchedCount, result.ModifiedCount)
	return err
}

func DeleteUser(userId string) error {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.M{"_id": id}
	collection := mongoClient.Database(db).Collection(dbColumn)
	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the users collection\n", result.DeletedCount)
	return err
}

func FindUserById(userId string) (User, error) {
	var result User
	filter := bson.M{"_id": userId}
	collection := mongoClient.Database(db).Collection(dbColumn)
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result, err
}

func FindUser(userId string) (User, error) {
	var result User
	filter := bson.D{{"_id", userId}}
	collection := mongoClient.Database(db).Collection(dbColumn)
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result, err
}

func FindAllUsers(userId string) ([]User, error) {
	var results []User
	filter := bson.D{{"_id", userId}}
	collection := mongoClient.Database(db).Collection(dbColumn)
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	err = cur.All(context.TODO(), &results)
	if err != nil {
		log.Fatal(err)
	}
	return results, err
}

func ListAll(userId string) ([]User, error) {
	var results []User
	collection := mongoClient.Database(db).Collection(dbColumn)
	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	err = cur.All(context.TODO(), &results)
	if err != nil {
		log.Fatal(err)
	}
	return results, err
}

func DeleteAll() error {
	collection := mongoClient.Database(db).Collection(dbColumn)
	result, err := collection.DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the users collection\n", result.DeletedCount)
	return err
}

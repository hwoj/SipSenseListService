package repositories

import (
	"context"
	"listservice/internal/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type ListRepository struct {
	MongoCollection *mongo.Collection
}

var MongoClient *mongo.Client

func ConnectDB() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	mongoURI := os.Getenv("MONGO_URI")

	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	MongoClient, err = mongo.Connect(context.Background(), opts)

	if err != nil {
		log.Fatal("Unable to connect to Mongo", err)
	}

	err = MongoClient.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("ping failed", err)
	}
}

func (listRepository ListRepository) CreateList(list *models.List) (interface{}, error) {
	result, err := listRepository.MongoCollection.InsertOne(context.Background(), listRepository)

	if err != nil {
		return nil, err
	}

	return result.InsertedID, err
}

func (listRepository ListRepository) GetListByID(ID string) (*models.List, error) {
	var list models.List

	err := listRepository.MongoCollection.FindOne(context.Background(), bson.D{{Key: "_id", Value: ID}}).Decode(&list)

	if err != nil {
		return nil, err
	}

	return &list, nil
}

func (listRepository ListRepository) DeleteList(ID string) (int64, error) {
	result, err := listRepository.MongoCollection.DeleteOne(context.Background(), bson.D{{Key: "_id", Value: ID}})

	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

func (listRepository ListRepository) MakeListPublic(ID string) (int64, error) {
	result, err := listRepository.MongoCollection.UpdateOne(context.Background(), bson.D{{Key: "_id", Value: ID}}, bson.D{{Key: "$set", Value: bson.D{{Key: "public", Value: true}}}})

	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, err
}

func (listRepository ListRepository) MakeListPrivate(ID string) (int64, error) {
	result, err := listRepository.MongoCollection.UpdateOne(context.Background(), bson.D{{Key: "_id", Value: ID}}, bson.D{{Key: "$set", Value: bson.D{{Key: "public", Value: false}}}})

	if err != nil {
		return 0, nil
	}

	return result.ModifiedCount, nil
}

func (listRepository ListRepository) AddDrinkToList(listID string, drinkID string) (int64, error) {
	result, err := listRepository.MongoCollection.UpdateOne(context.Background(), bson.M{"_id": listID}, bson.M{"$push": bson.M{"drinks": drinkID}})

	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, err
}

func (listRepository ListRepository) RemoveDrinkFromList(listID string, drinkID string) (int64, error) {
	result, err := listRepository.MongoCollection.UpdateOne(context.Background(), bson.M{"_id": listID}, bson.M{"$pull": bson.M{"drinks": drinkID}})

	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, err
}

func (listRepository ListRepository) AddCollaborator(listID string, userID string) (int64, error) {
	result, err := listRepository.MongoCollection.UpdateOne(context.Background(), bson.M{"_id": listID}, bson.M{"$push": bson.M{"collaborators": userID}})

	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, err
}

func (listRepository ListRepository) RemoveCollaborator(listID string, userID string) (int64, error) {
	result, err := listRepository.MongoCollection.UpdateOne(context.Background(), bson.M{"_id": listID}, bson.M{"$pull": bson.M{"drinks": userID}})

	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, err
}

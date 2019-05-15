package persistence

import (
	"context"

	"github.com/shauera/messages/model"
	"github.com/shauera/messages/utils"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"

	log "github.com/sirupsen/logrus"
	config "github.com/spf13/viper"
)

type MongoRepository struct {
	ctx          context.Context
	client       *mongo.Client
	databaseName string
}

func NewMongoRepository(ctx context.Context) *MongoRepository {
	mongoConnectionString := `mongodb://` + config.GetString("database.server")
	username := config.GetString("database.username")
	password := config.GetString("database.password")

	clientOptions := options.Client().SetAuth(options.Credential{Username: username, Password: password})
	client, err := mongo.Connect(ctx, mongoConnectionString, clientOptions)
	if err != nil {
		log.WithError(err).Warning("Could not connect to database")
	}

	return &MongoRepository{
		ctx:          ctx,
		client:       client,
		databaseName: config.GetString("database.dbname"),
	}
}

func (mr *MongoRepository) CreateMessage(message model.MessageRequest) (*model.MessageResponse, error) {
	createMessage := model.MessageResponse{
		Author: message.Author,
		Content: message.Content,
		CreatedAt: message.CreatedAt,
		Palindrome: utils.IsPalindrome(message.Content),
	}

	collection := mr.client.Database(mr.databaseName).Collection("messages")
	result, err := collection.InsertOne(mr.ctx, createMessage)
	if err != nil {
		return nil, err
	}

	createMessage.ID = result.InsertedID.(primitive.ObjectID).Hex()
	
	return &createMessage, nil
}

func (mr *MongoRepository) ListMessages() (model.MessageResponses, error) {
	collection := mr.client.Database(mr.databaseName).Collection("messages")

	cursor, err := collection.Find(mr.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(mr.ctx)

	var MessageResponses model.MessageResponses
	for cursor.Next(mr.ctx) {
		var MessageResponse model.MessageResponse
		cursor.Decode(&MessageResponse)
		MessageResponses = append(MessageResponses, MessageResponse)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return MessageResponses, nil
}

func (mr *MongoRepository) FindMessageById(id string) (*model.MessageResponse, error) {
	collection := mr.client.Database(mr.databaseName).Collection("messages")

	messageID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var message model.MessageResponse
	err = collection.FindOne(mr.ctx, model.MessageResponse{ID: messageID}).Decode(&message)
	if err != nil && err.Error() == "mongo: no documents in result" {
		return nil, ErrorNotFound
	}

	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (mr *MongoRepository) UpdateMessageById(id string, message model.MessageRequest) (*model.MessageResponse, error) {
	collection := mr.client.Database(mr.databaseName).Collection("messages")

	messageID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := model.MessageResponse{ID: messageID}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "author", Value: message.Author},
			{Key: "content", Value: message.Content},
			{Key: "createdAt", Value: message.CreatedAt},
			{Key: "palindrome", Value: utils.IsPalindrome(message.Content)},
		}},
	}
	updateOptions := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedMessage model.MessageResponse
	err = collection.FindOneAndUpdate(mr.ctx, filter, update, updateOptions).Decode(&updatedMessage)
	if err != nil && err.Error() == "mongo: no documents in result" {
		return nil, ErrorNotFound
	}

	if err != nil {
		return nil, err
	}

	return &updatedMessage, nil
}

func (mr *MongoRepository) DeleteMessageById(id string) (error) {
	collection := mr.client.Database(mr.databaseName).Collection("messages")

	messageID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := collection.DeleteOne(mr.ctx, model.MessageResponse{ID: messageID})
	if err == nil && result.DeletedCount == 0 {
		return ErrorNotFound
	}

	return err
}
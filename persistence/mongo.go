package persistence

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/shauera/messages/application"
	"github.com/shauera/messages/model"
	"github.com/shauera/messages/utils"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"

	log "github.com/sirupsen/logrus"
	config "github.com/spf13/viper"
)

//MongoRepository - mongo collection (database) for persisting message documents
type MongoRepository struct {
	client       *mongo.Client
	databaseName string
}

//NewMongoRepository - initialize and return a new MongoRepository
func NewMongoRepository(ctx context.Context, healthSystemID string) (*MongoRepository, error) {
	mongoConnectionString := `mongodb://` + config.GetString("database.server")
	username := config.GetString("database.username")
	password := config.GetString("database.password")

	repositoryContext, cancel := getRepositoryContext(ctx)
	defer cancel()

	clientOptions := options.Client().SetAuth(options.Credential{Username: username, Password: password})
	client, err := mongo.Connect(repositoryContext, mongoConnectionString, clientOptions)
	if err != nil {
		log.WithError(err).Debug("Could not connect to database")
		return nil, errors.Wrap(err, "Could not connect to database")
	}
	err = client.Ping(repositoryContext, nil)
	if err != nil {
		log.WithError(err).WithField("server", mongoConnectionString).Debug("Could not ping database")
		return nil, errors.Wrap(err, "Could not ping database")
	}

	go func() {
		<-ctx.Done()
		log.Debug("Closing mongodb connection")
		if client.Disconnect(ctx) != nil {
			log.Debug("Failed to close mongodb connection")
		}
		log.Debug("Mongodb connection closed")
	}()

	application.AddHealthMonitor(newMongoHealthMonitor(healthSystemID, client))

	return &MongoRepository{
		client:       client,
		databaseName: config.GetString("database.dbname"),
	}, nil
}

//CreateMessage - adds a new message record into repository
func (mr *MongoRepository) CreateMessage(ctx context.Context, message model.MessageRequest) (*model.MessageResponse, error) {
	repositoryContext, cancel := getRepositoryContext(ctx)
	defer cancel()

	createMessage := model.MessageResponse{
		ID:         primitive.NewObjectID(),
		Author:     updateString(nil, message.Author),
		Content:    updateString(nil, message.Content),
		CreatedAt:  (*model.MessageTime)(updateTime(nil, (*time.Time)(message.CreatedAt))),
		Palindrome: utils.IsPalindrome(*message.Content),
	}

	collection := mr.client.Database(mr.databaseName).Collection("messages")
	result, err := collection.InsertOne(repositoryContext, createMessage)
	if err != nil {
		return nil, err
	}

	createMessage.ID = result.InsertedID.(primitive.ObjectID).Hex()

	return &createMessage, nil
}

//UpdateMessageByID - updates an existing message record
//An error will be returned if the given id does not exist
func (mr *MongoRepository) UpdateMessageByID(ctx context.Context, id string, updateMessage model.MessageRequest) (*model.MessageResponse, error) {
	repositoryContext, cancel := getRepositoryContext(ctx)
	defer cancel()

	oldMessage, err := mr.FindMessageByID(repositoryContext, id)
	if err != nil {
		return nil, err
	}

	collection := mr.client.Database(mr.databaseName).Collection("messages")

	messageID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := model.MessageResponse{ID: messageID}

	updatedAuthor := updateString(oldMessage.Author, updateMessage.Author)
	updatedCreatedAt := (*model.MessageTime)(updateTime((*time.Time)(oldMessage.CreatedAt), (*time.Time)(updateMessage.CreatedAt)))
	upadtedContent := updateString(oldMessage.Content, updateMessage.Content)
	var updatedPalindrome bool
	if updateMessage.Content != nil && *upadtedContent != *oldMessage.Content {
		// message content got a new value, calculating new palindrome state
		updatedPalindrome = utils.IsPalindrome(*updateMessage.Content)
	}

	update := bson.D{
		//{Key: "$set", Value: bson.D{
		{Key: op(updatedAuthor), Value: bson.D{
			{Key: "author", Value: updatedAuthor},
		}},
		{Key: op(upadtedContent), Value: bson.D{
			{Key: "content", Value: upadtedContent},
		}},
		{Key: op(updatedCreatedAt), Value: bson.D{
			{Key: "createdAt", Value: updatedCreatedAt},
		}},
		{Key: "$set", Value: bson.D{
			{Key: "palindrome", Value: updatedPalindrome},
		}},
	}

	updateOptions := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedMessage model.MessageResponse
	err = collection.FindOneAndUpdate(repositoryContext, filter, update, updateOptions).Decode(&updatedMessage)
	if err != nil && err.Error() == "mongo: no documents in result" {
		return nil, ErrorNotFound
	}

	if err != nil {
		return nil, err
	}

	return &updatedMessage, nil
}

//ListMessages - returns all message records in the repository
func (mr *MongoRepository) ListMessages(ctx context.Context) (model.MessageResponses, error) {
	repositoryContext, cancel := getRepositoryContext(ctx)
	defer cancel()

	collection := mr.client.Database(mr.databaseName).Collection("messages")

	cursor, err := collection.Find(repositoryContext, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(repositoryContext)

	var MessageResponses model.MessageResponses
	for cursor.Next(repositoryContext) {
		var MessageResponse model.MessageResponse
		cursor.Decode(&MessageResponse)
		MessageResponses = append(MessageResponses, MessageResponse)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return MessageResponses, nil
}

//FindMessageByID - returns an existing message record
//An error will be returned if the given id does not exist
func (mr *MongoRepository) FindMessageByID(ctx context.Context, id string) (*model.MessageResponse, error) {
	repositoryContext, cancel := getRepositoryContext(ctx)
	defer cancel()

	collection := mr.client.Database(mr.databaseName).Collection("messages")

	messageID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var messageResponse model.MessageResponse
	err = collection.FindOne(repositoryContext, model.MessageResponse{ID: messageID}).Decode(&messageResponse)
	if err != nil && err.Error() == "mongo: no documents in result" {
		return nil, ErrorNotFound
	}

	if err != nil {
		return nil, err
	}

	return &messageResponse, nil
}

//DeleteMessageByID - removes an existing message record from the repository
//An error will be returned if the given id does not exist
func (mr *MongoRepository) DeleteMessageByID(ctx context.Context, id string) error {
	repositoryContext, cancel := getRepositoryContext(ctx)
	defer cancel()

	collection := mr.client.Database(mr.databaseName).Collection("messages")

	messageID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := collection.DeleteOne(repositoryContext, model.MessageResponse{ID: messageID})
	if err == nil && result.DeletedCount == 0 {
		return ErrorNotFound
	}

	return err
}

func op(value interface{}) string {
	if !utils.IsNilValue(value) {
		return "$set"
	}

	return "$unset"
}

func getRepositoryContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, config.GetDuration("database.timeout"))
}

//--------------------------- Health Monitor ------------------------------------
type mongoHealthMonitor struct {
	healthSystemID string
	client         *mongo.Client
}

func newMongoHealthMonitor(healthSystemID string, client *mongo.Client) application.HealthMonitor {
	return mongoHealthMonitor{
		healthSystemID: healthSystemID,
		client:         client,
	}
}

func (mhm mongoHealthMonitor) IsHealthy(ctx context.Context) bool {
	repositoryContext, cancel := getRepositoryContext(ctx)
	defer cancel()

	err := mhm.client.Ping(repositoryContext, nil)
	if err != nil {
		log.WithError(err).WithField("systemID", mhm.healthSystemID).Debug("Could not ping database")
		return false
	}

	return true
}

func (mhm mongoHealthMonitor) GetSystemID() string {
	return mhm.healthSystemID
}

package persistence

import (
	"context"

	"github.com/shauera/messages/model"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"

	log "github.com/sirupsen/logrus"
	config "github.com/spf13/viper"
)

type MongoRepository struct {
	ctx context.Context
	client *mongo.Client
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
		ctx: ctx,
		client: client,
		databaseName: config.GetString("database.dbname"),
	}
}

func (mr *MongoRepository) CreatePerson(person model.Person) (*string, error) {
	collection := mr.client.Database(mr.databaseName).Collection("people")
	result, err := collection.InsertOne(mr.ctx, person)
	if err != nil {
		return nil, err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()
	return &id, nil
}

func (mr *MongoRepository) ListPersons() (model.Persons, error) {
	collection := mr.client.Database(mr.databaseName).Collection("people")

	cursor, err := collection.Find(mr.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(mr.ctx)

	var people model.Persons
	for cursor.Next(mr.ctx) {
		var person model.Person
		cursor.Decode(&person)
		people = append(people, person)
	}

	if err := cursor.Err(); err != nil {
		return nil ,err
	}

	return people, nil
}

func (mr *MongoRepository) FindPersonById(id string) (*model.Person, error) {
	collection := mr.client.Database(mr.databaseName).Collection("people")

	personID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	
	var person model.Person
	err = collection.FindOne(mr.ctx, model.Person{ID: personID}).Decode(&person)
	if err != nil {
		return nil, err
	}

	return &person, nil
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	config "github.com/spf13/viper"
)

// utcJSONFormatter is a logrus JSON Formatter wrapper that forces the event time to be UTC on format.
type utcJSONFormatter struct {
	fmt *log.JSONFormatter
}

// Format the log entry by forcing the event time to be UTC and delegating to the wrapped Formatter.
func (u utcJSONFormatter) Format(e *log.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return u.fmt.Format(e)
}

func init() {
	// ---------------- Configuration -----------------------------
	// Source
	config.SetConfigName("config")
	config.AddConfigPath("/etc/messages")
	config.AddConfigPath(".")

	err := config.ReadInConfig() // Find and read the config file
	if err != nil {
		log.WithError(err).Fatal("Could not read config file")
	}

	// Defaults
	config.SetDefault(
		"logging", map[string]interface{}{
			"level": "info",
		},
	)

	config.SetDefault(
		"service", map[string]interface{}{
			"port": "8090",
		},
	)

	config.SetDefault(
		"mongo", map[string]interface{}{
			"server":   "mongo:27017",
			"username": "root",
			"password": "example",
			"timeout":  10,
		},
	)

	// ---------------- Configure logging -------------------------

	// Set minimal logging level
	logLevel, err := log.ParseLevel(config.GetString("logging.level"))
	if err != nil {
		logLevel = log.InfoLevel
	}
	log.SetLevel(logLevel)

	// Log file, function and line number
	log.SetReportCaller(true)

	// Log as JSON instead of the default ASCII formatter enforcing UTC timezone
	formatter := utcJSONFormatter{fmt: new(log.JSONFormatter)}

	// Format log to show short file and function names
	gitPath := "github.com/shauera/messages"
	repoPath := fmt.Sprintf("%s/src/"+gitPath, os.Getenv("GOPATH"))
	formatter.fmt.CallerPrettyfier = func(f *runtime.Frame) (string, string) {
		fileName := strings.Replace(f.File, repoPath, "", -1)
		functionName := strings.Replace(f.Function, gitPath, "", -1)
		return fmt.Sprintf("%s()", functionName), fmt.Sprintf("%s:%d", fileName, f.Line)
	}
	log.SetFormatter(formatter)
}

type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

type errorResponse struct {
	Message   string `json:"message" bson:"message"`
}

func CreatePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var person Person
	_ = json.NewDecoder(request.Body).Decode(&person)
	collection := client.Database("thepolyglotdeveloper").Collection("people")
	ctx, cancel := context.WithTimeout(request.Context(), config.GetDuration("mongo.timeout")*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)		
		json.NewEncoder(response).Encode(errorResponse{Message: err.Error()})
		log.WithError(err).Debug("Could not get person")

		return
	}

	json.NewEncoder(response).Encode(result)
}

func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var people []Person
	collection := client.Database("thepolyglotdeveloper").Collection("people")
	ctx, cancel := context.WithTimeout(request.Context(), config.GetDuration("mongo.timeout")*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(errorResponse{Message: err.Error()})
		log.WithError(err).Debug("Could not get people")
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(errorResponse{Message: err.Error()})
		log.WithError(err).Debug("Could not get people")

		return
	}
	json.NewEncoder(response).Encode(people)
}

func GetPersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person Person
	collection := client.Database("thepolyglotdeveloper").Collection("people")
	ctx, cancel := context.WithTimeout(request.Context(), config.GetDuration("mongo.timeout")*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, Person{ID: id}).Decode(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(errorResponse{Message: err.Error()})
		log.WithError(err).Debug("Could not get person")

		return
	}
	json.NewEncoder(response).Encode(person)
}

var client *mongo.Client

func setupMux() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/person", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/person/{id}", GetPersonEndpoint).Methods("GET")
	return router
}

func main() {
	fmt.Println("Starting the application...")

	ctx, cancel := context.WithTimeout(context.Background(), config.GetDuration("mongo.timeout")*time.Second)
	defer cancel()

	mongoConnectionString := `mongodb://` + config.GetString("mongo.server")
	username := config.GetString("mongo.username")
	password := config.GetString("mongo.password")

	clientOptions := options.Client().SetAuth(options.Credential{Username: username, Password: password})
	var err error
	client, err = mongo.Connect(ctx, mongoConnectionString, clientOptions)
	if err != nil {
		log.WithError(err).Fatal("Could not connect to database - aborting")
	}

	router := setupMux()

	bindPort := ":" + config.GetString("service.port")
	log.Fatal(http.ListenAndServe(bindPort, router))
}

package rest

import(
	"context"
	"encoding/json"
	"net/http"

	"github.com/shauera/messages/model"
	modelCommon "github.com/shauera/messages/model/common"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

type PersonRepository interface{
	FindPersonById(id string) (*model.Person, error)
	CreatePerson(person model.Person) (*string, error)
    ListPersons() (model.Persons, error)
}

type PersonController struct {
	ctx context.Context
	repository PersonRepository
}

func (pc *PersonController) CreatePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var person model.Person
	_ = json.NewDecoder(request.Body).Decode(&person)

	personID, err := pc.repository.CreatePerson(person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)		
		json.NewEncoder(response).Encode(modelCommon.ErrorResponse{Message: err.Error()})
		log.WithError(err).Debug("Could not create person")

		return
	}

	json.NewEncoder(response).Encode(personID)
}

func (pc *PersonController) GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) {
	// swagger:operation GET /hello/{name} hello Hello
	//
	// Returns a simple Hello message
	// ---
	// consumes:
	// - text/plain
	// produces:
	// - text/plain
	// parameters:
	// - name: name
	//   in: path
	//   description: Name to be returned.
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: The hello message
	//     type: string
	response.Header().Set("content-type", "application/json")

	persons, err := pc.repository.ListPersons()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)		
		json.NewEncoder(response).Encode(modelCommon.ErrorResponse{Message: err.Error()})
		log.WithError(err).Debug("Could not get list of persons")

		return
	}
	json.NewEncoder(response).Encode(persons)
}

func (pc *PersonController) GetPersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	person, err := pc.repository.FindPersonById(params["id"])
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(modelCommon.ErrorResponse{Message: err.Error()})
		log.WithError(err).Debug("Could not get person")

		return
	}
	json.NewEncoder(response).Encode(person)
}
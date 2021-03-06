package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/shauera/messages/model"
	modelCommon "github.com/shauera/messages/model"
	"github.com/shauera/messages/persistence"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

// MessageRepository - repository abstraction to be implemented by persisters
type MessageRepository interface {
	FindMessageByID(ctx context.Context, id string) (*model.MessageResponse, error)
	CreateMessage(ctx context.Context, message model.MessageRequest) (*model.MessageResponse, error)
	ListMessages(ctx context.Context) (model.MessageResponses, error)
	DeleteMessageByID(ctx context.Context, id string) error
	UpdateMessageByID(ctx context.Context, id string, message model.MessageRequest) (*model.MessageResponse, error)
}

// MessageController - handles message resource endpoints
type MessageController struct {
	repository MessageRepository
}

//NewMessageController - return a new message controller setup with a designated message repository
func NewMessageController(messageRepository MessageRepository) MessageController {
	return MessageController{
		repository: messageRepository,
	}
}

//PublishEndpoints - implementation of ServiceController
func (mc MessageController) PublishEndpoints(router *mux.Router) {
	router.HandleFunc("/messages", mc.CreateMessage).Methods("POST")
	router.HandleFunc("/messages", mc.ListMessages).Methods("GET")
	router.HandleFunc("/messages/{id}", mc.GetMessageByID).Methods("GET")
	router.HandleFunc("/messages/{id}", mc.UpdateMessageByID).Methods("PUT")
	router.HandleFunc("/messages/{id}", mc.DeleteMessageByID).Methods("DELETE")
}

//------------------------------- Create -----------------------------------------

// CreateMessage - creates a new message
func (mc *MessageController) CreateMessage(response http.ResponseWriter, request *http.Request) {
	// swagger:operation POST /messages messages createMessage
	//
	// Creates a new message
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: messageRequest
	//   in: body
	//   description: message to be created.
	//   required: false
	//   type: MessageRequest
	//   schema:
	//     "$ref": "#/definitions/MessageRequest"
	// responses:
	//   '200':
	//     description: OK
	//     type: string
	//   '400':
	//     description: Bad Request
	//   '500':
	//     description: Internal Server Error

	newMessage, err := validateRequest(response, request)
	if err != nil {
		return
	}

	messageID, err := mc.repository.CreateMessage(request.Context(), *newMessage)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(modelCommon.ErrorResponse{Message: err.Error()})
		log.WithError(err).Debug("Could not create message")
		return
	}

	json.NewEncoder(response).Encode(messageID)
}

//------------------------------- Gel All ----------------------------------------
// TODO - add filtering
//      - add pagination

// ListMessages - retrieves a list of all available messages
func (mc *MessageController) ListMessages(response http.ResponseWriter, request *http.Request) {
	// swagger:operation GET /messages messages listMessages
	//
	// Returns a list of all available messages
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: OK
	//     schema:
	//       "$ref": "#/definitions/MessageResponses"
	//   '500':
	//     description: Internal Server Error
	response.Header().Set("content-type", "application/json")

	messages, err := mc.repository.ListMessages(request.Context())
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(modelCommon.ErrorResponse{Message: err.Error()})
		log.WithError(err).Debug("Could not get list of messages")

		return
	}
	json.NewEncoder(response).Encode(messages)
}

//------------------------------- Get --------------------------------------------

// GetMessageByID - retrieves a single message by id
func (mc *MessageController) GetMessageByID(response http.ResponseWriter, request *http.Request) {
	// swagger:operation GET /messages/{id} messages listMessage
	//
	// Returns a message by id
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   description: id of message to be returned.
	//   required: true
	//   type: string
	// responses:
	//   '200':
	//     description: OK
	//     schema:
	//       "$ref": "#/definitions/MessageResponse"
	//   '404':
	//     description: Not Found
	//   '500':
	//     description: Internal Server Error
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	message, err := mc.repository.FindMessageByID(request.Context(), params["id"])
	if err != nil {
		if err == persistence.ErrorNotFound {
			response.WriteHeader(http.StatusNotFound)
		} else {
			response.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(response).Encode(modelCommon.ErrorResponse{Message: err.Error()})
			log.WithError(err).Debug("Could not get message")
		}
		return
	}
	json.NewEncoder(response).Encode(message)
}

//------------------------------- Update -----------------------------------------

// UpdateMessageByID - updates an existing message
func (mc *MessageController) UpdateMessageByID(response http.ResponseWriter, request *http.Request) {
	// swagger:operation PUT /messages/{id} messages updateMessage
	//
	// Updates message by id
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   description: id of message to be updated.
	//   required: true
	//   type: string
	// - name: messageRequest
	//   in: body
	//   description: message to be updated.
	//   required: false
	//   type: MessageRequest
	//   schema:
	//     "$ref": "#/definitions/MessageRequest"
	// responses:
	//   '200':
	//     description: OK
	//     schema:
	//       "$ref": "#/definitions/MessageResponse"
	//   '404':
	//     description: Not Found
	//   '500':
	//     description: Internal Server Error

	updatedMessage, err := validateRequest(response, request)
	if err != nil {
		return
	}

	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	message, err := mc.repository.UpdateMessageByID(request.Context(), params["id"], *updatedMessage)
	if err != nil {
		if err == persistence.ErrorNotFound {
			response.WriteHeader(http.StatusNotFound)
		} else {
			response.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(response).Encode(modelCommon.ErrorResponse{Message: err.Error()})
			log.WithError(err).Debug("Could not update message")
		}
		return
	}
	json.NewEncoder(response).Encode(message)
}

//------------------------------- Delete -----------------------------------------

// DeleteMessageByID - deletes an existing message
func (mc *MessageController) DeleteMessageByID(response http.ResponseWriter, request *http.Request) {
	// swagger:operation DELETE /messages/{id} messages deleteMessage
	//
	// Delete a message by id
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// parameters:
	// - name: id
	//   in: path
	//   description: id of message to be deleted.
	//   required: true
	//   type: string
	// responses:
	//   '204':
	//     description: No Content
	//   '404':
	//     description: Not Found
	//   '500':
	//     description: Internal Server Error
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	err := mc.repository.DeleteMessageByID(request.Context(), params["id"])
	if err != nil {
		if err == persistence.ErrorNotFound {
			response.WriteHeader(http.StatusNotFound)
		} else {
			response.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(response).Encode(modelCommon.ErrorResponse{Message: err.Error()})
			log.WithError(err).Debug("Could not delete message")
		}
		return
	}
	response.WriteHeader(http.StatusNoContent)
}

//------------------------------- Validation -------------------------------------

func validateRequest(response http.ResponseWriter, request *http.Request) (*model.MessageRequest, error) {
	response.Header().Set("content-type", "application/json")

	var newMessage model.MessageRequest
	err := json.NewDecoder(request.Body).Decode(&newMessage)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		responseErr := errors.Wrap(err, "Could not decode request body")
		json.NewEncoder(response).Encode(modelCommon.ErrorResponse{Message: responseErr.Error()})
		log.WithError(err).Debug("Could not decode request body")
		return nil, errors.New("validation failed")
	}

	validationErrorsResponse := newMessage.Validate()
	if len(validationErrorsResponse.Messages) != 0 {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(validationErrorsResponse)
		log.Debug("Validation of message request failed")
		return nil, errors.New("validation failed")
	}

	return &newMessage, nil
}

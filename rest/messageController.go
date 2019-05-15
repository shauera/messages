package rest

import (
	"context"
	"encoding/json"
	"github.com/shauera/messages/persistence"
	"net/http"

	"github.com/shauera/messages/model"
	modelCommon "github.com/shauera/messages/model/common"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

// MessageRepository - repository abstraction to be implemented by persisters
type MessageRepository interface {
	FindMessageById(id string) (*model.MessageResponse, error)
	CreateMessage(message model.MessageRequest) (*model.MessageResponse, error)
	ListMessages() (model.MessageResponses, error)
	DeleteMessageById(id string) error
	UpdateMessageById(id string, message model.MessageRequest) (*model.MessageResponse, error)
}

// MessageController - handles message resource endpoints
type MessageController struct {
	ctx        context.Context
	repository MessageRepository
}

//------------------------------- handlers ---------------------------------------

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
	response.Header().Set("content-type", "application/json")
	var newMessage model.MessageRequest
	err := json.NewDecoder(request.Body).Decode(&newMessage)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(modelCommon.ErrorResponse{Message: err.Error()})
		log.WithError(err).Debug("Could not decode request body")
		return
	}

	messageID, err := mc.repository.CreateMessage(newMessage)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(modelCommon.ErrorResponse{Message: err.Error()})
		log.WithError(err).Debug("Could not create message")
		return
	}

	json.NewEncoder(response).Encode(messageID)
}

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

	messages, err := mc.repository.ListMessages()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(modelCommon.ErrorResponse{Message: err.Error()})
		log.WithError(err).Debug("Could not get list of messages")

		return
	}
	json.NewEncoder(response).Encode(messages)
}

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
	message, err := mc.repository.FindMessageById(params["id"])
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
	//   description: id of message to be returned.
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
	var updatedMessage model.MessageRequest
	err := json.NewDecoder(request.Body).Decode(&updatedMessage)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(modelCommon.ErrorResponse{Message: err.Error()})
		log.WithError(err).Debug("Could not decode request body")
		return
	}

	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	message, err := mc.repository.UpdateMessageById(params["id"], updatedMessage)
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
	err := mc.repository.DeleteMessageById(params["id"])
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

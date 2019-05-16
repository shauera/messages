package persistence

import (
	"strconv"
	"sync/atomic"
	"time"

	"github.com/shauera/messages/model"
	"github.com/shauera/messages/utils"
)

type MemoryRepository struct {
	messageIDCounter int64
	messagesStorage  map[string]model.MessageResponse
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		messagesStorage: make(map[string]model.MessageResponse),
	}
}

func (mr *MemoryRepository) CreateMessage(newMessage model.MessageRequest) (*model.MessageResponse, error) {
	id := strconv.FormatInt(atomic.AddInt64(&mr.messageIDCounter, 1), 10)

	messageResponse := mr.storeMessage(id, model.MessageResponse{}, newMessage)
	return messageResponse, nil
}

func (mr *MemoryRepository) UpdateMessageById(id string, updateMessage model.MessageRequest) (*model.MessageResponse, error) {

	if oldMessage, ok := mr.messagesStorage[id]; ok {
		return mr.storeMessage(id, oldMessage, updateMessage), nil
	}

	return nil, ErrorNotFound
}

func (mr *MemoryRepository) ListMessages() (model.MessageResponses, error) {

	messageResponses := make(model.MessageResponses, 0, len(mr.messagesStorage))

	for _, value := range mr.messagesStorage {
		messageResponses = append(messageResponses, value)
	}

	if len(messageResponses) == 0 {
		return nil, nil
	}

	return messageResponses, nil
}

func (mr *MemoryRepository) FindMessageById(id string) (*model.MessageResponse, error) {
	if messageResponse, ok := mr.messagesStorage[id]; ok {
		return &messageResponse, nil
	}

	return nil, ErrorNotFound
}

func (mr *MemoryRepository) DeleteMessageById(id string) error {
	if _, ok := mr.messagesStorage[id]; ok {
		delete(mr.messagesStorage, id)
		return nil
	}

	return ErrorNotFound
}

func updateString(oldValue, newValue *string) *string {
	if newValue != nil {
		if *newValue == "" { // update is explicitly removing the field
			return nil
		}
		return newValue // update is explicitly setting the field to a new value
	}
	return oldValue // update did not explicitly set this field
}

func updateTime(oldValue, newValue *model.MessageTime) *model.MessageTime {
	if newValue != nil {
		tmpTime := time.Time(*newValue)
		if time.Time.IsZero(tmpTime) { // update is explicitly removing the field
			return nil
		}
		return newValue // update is explicitly setting the field to a new value
	}
	return oldValue // update did not explicitly set this field
}

func (mr *MemoryRepository) storeMessage(id string, oldMessage model.MessageResponse, updateMessage model.MessageRequest) *model.MessageResponse {
	newMessageResponse := model.MessageResponse{
		ID:        id,
		Author:    updateString(oldMessage.Author, updateMessage.Author),
		Content:   updateString(oldMessage.Content, updateMessage.Content),
		CreatedAt: updateTime(oldMessage.CreatedAt, updateMessage.CreatedAt),
	}

	if oldMessage.Content == nil ||
		updateMessage.Content != nil && *newMessageResponse.Content != *oldMessage.Content {
		// message content got a new value, calculating new palindrome state
		newMessageResponse.Palindrome = utils.IsPalindrome(*newMessageResponse.Content)
	}

	mr.messagesStorage[id] = newMessageResponse

	return &newMessageResponse
}

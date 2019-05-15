package persistence

import (
	"strconv"
	"sync/atomic"

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

	messageResponse := mr.storeMessage(id, newMessage)
	return messageResponse, nil
}

func (mr *MemoryRepository) UpdateMessageById(id string, updateMessage model.MessageRequest) (*model.MessageResponse, error) {

	if _, ok := mr.messagesStorage[id]; ok {
		return nil, ErrorNotFound
	}

	return mr.storeMessage(id, updateMessage), nil
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

func (mr *MemoryRepository) DeleteMessageById(id string) (error) {
	if _, ok := mr.messagesStorage[id]; ok {
		delete(mr.messagesStorage, id)
		return nil
	}

	return ErrorNotFound	
}

func (mr *MemoryRepository) storeMessage(id string, updateMessage model.MessageRequest) (*model.MessageResponse) {
	newMessageResponse := model.MessageResponse {
		ID: id,
		Author: updateMessage.Author,
		Content: updateMessage.Content,
		CreatedAt: updateMessage.CreatedAt,
		Palindrome: utils.IsPalindrome(updateMessage.Content),
	}
	mr.messagesStorage[id] = newMessageResponse

	return &newMessageResponse
}

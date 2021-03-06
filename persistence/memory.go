package persistence

import (
	"context"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/shauera/messages/model"
	"github.com/shauera/messages/utils"
)

//MemoryRepository - in memory repository for use with demo, mocking out real database and tests
//ALL RECORDS WILL BE DELETED ONCE THE INSTANCE IS RESTARTED!
type MemoryRepository struct {
	messageIDCounter int64
	messagesStorage  map[string]model.MessageResponse
}

//NewMemoryRepository - initialize and return a new MemoryRepository
func NewMemoryRepository() (*MemoryRepository, error) {
	return &MemoryRepository{
		messagesStorage: make(map[string]model.MessageResponse),
	}, nil
}

//CreateMessage - adds a new message record into repository
func (mr *MemoryRepository) CreateMessage(ctx context.Context, newMessage model.MessageRequest) (*model.MessageResponse, error) {
	id := strconv.FormatInt(atomic.AddInt64(&mr.messageIDCounter, 1), 10)

	messageResponse := mr.storeMessage(ctx, id, model.MessageResponse{}, newMessage)
	return messageResponse, nil
}

//UpdateMessageByID - updates an existing message record
//An error will be returned if the given id does not exist 
func (mr *MemoryRepository) UpdateMessageByID(ctx context.Context, id string, updateMessage model.MessageRequest) (*model.MessageResponse, error) {

	if oldMessage, ok := mr.messagesStorage[id]; ok {
		return mr.storeMessage(ctx, id, oldMessage, updateMessage), nil
	}

	return nil, ErrorNotFound
}

//ListMessages - returns all message records in the repository
func (mr *MemoryRepository) ListMessages(ctx context.Context) (model.MessageResponses, error) {

	messageResponses := make(model.MessageResponses, 0, len(mr.messagesStorage))

	for _, value := range mr.messagesStorage {
		messageResponses = append(messageResponses, value)
	}

	if len(messageResponses) == 0 {
		return nil, nil
	}

	return messageResponses, nil
}

//FindMessageByID - returns an existing message record
//An error will be returned if the given id does not exist  
func (mr *MemoryRepository) FindMessageByID(ctx context.Context, id string) (*model.MessageResponse, error) {
	if messageResponse, ok := mr.messagesStorage[id]; ok {
		return &messageResponse, nil
	}

	return nil, ErrorNotFound
}

//DeleteMessageByID - removes an existing message record from the repository
//An error will be returned if the given id does not exist 
func (mr *MemoryRepository) DeleteMessageByID(ctx context.Context, id string) error {
	if _, ok := mr.messagesStorage[id]; ok {
		delete(mr.messagesStorage, id)
		return nil
	}

	return ErrorNotFound
}

func (mr *MemoryRepository) storeMessage(ctx context.Context, id string, oldMessage model.MessageResponse, updateMessage model.MessageRequest) *model.MessageResponse {
	newMessageResponse := model.MessageResponse{
		ID:        id,
		Author:    updateString(oldMessage.Author, updateMessage.Author),
		Content:   updateString(oldMessage.Content, updateMessage.Content),
		CreatedAt: (*model.MessageTime)(updateTime((*time.Time)(oldMessage.CreatedAt), (*time.Time)(updateMessage.CreatedAt))),
	}

	if oldMessage.Content == nil ||
		updateMessage.Content != nil && *newMessageResponse.Content != *oldMessage.Content {
		// message content got a new value, calculating new palindrome state
		newMessageResponse.Palindrome = utils.IsPalindrome(*newMessageResponse.Content)
	}

	mr.messagesStorage[id] = newMessageResponse

	return &newMessageResponse
}

//GetMessagesStorage - allows direct manipualtion of the storage to facilitate testing
func (mr *MemoryRepository) GetMessagesStorage() map[string]model.MessageResponse {
	return mr.messagesStorage
}
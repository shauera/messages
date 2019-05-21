package rest

import (
	"fmt"
	"time"

	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shauera/messages/model"

	"github.com/stretchr/testify/assert"

	"github.com/shauera/messages/persistence"
)

func getNewString(str string) *string {
	return &str
}

func getNewMessageTime(t time.Time) *model.MessageTime {
	mt := model.MessageTime(t)
	return &mt
}

//------------------------------- Gel All ----------------------------------------
func Test_Get_All(t *testing.T) {
	testCases := []struct {
		name    string
		preload func(memoryRepository *persistence.MemoryRepository)
		checker func(t *testing.T, response *httptest.ResponseRecorder)
	}{
		{
			name:    "Success path - empty repository",
			preload: func(memoryRepository *persistence.MemoryRepository) {},
			checker: func(t *testing.T, response *httptest.ResponseRecorder) {
				assert.Equal(t, "null\n", response.Body.String())
				assert.Equal(t, http.StatusOK, response.Code)
			},
		},
		{
			name: "Success path - 2 records in repository",
			preload: func(memoryRepository *persistence.MemoryRepository) {
				messagesStorage := memoryRepository.GetMessagesStorage()
				messagesStorage["8"] = model.MessageResponse{
					ID:         "8",
					Content:    getNewString("Test Message 1"),
					Author:     getNewString("test author 1"),
					CreatedAt:  getNewMessageTime(time.Date(2016, time.August, 15, 0, 0, 0, 0, time.UTC)),
					Palindrome: false,
				}
				messagesStorage["10"] = model.MessageResponse{
					ID:         "10",
					Content:    getNewString("Test Message 2"),
					Author:     getNewString("test author 2"),
					CreatedAt:  getNewMessageTime(time.Date(2017, time.August, 15, 0, 0, 0, 0, time.UTC)),
					Palindrome: false,
				}
			},
			checker: func(t *testing.T, response *httptest.ResponseRecorder) {
				assert.Contains(t, response.Body.String(), "{\"id\":\"8\",\"content\":\"Test Message 1\",\"author\":\"test author 1\",\"createdAt\":\"2016-08-15T00:00:00Z\",\"palindrome\":false}")
				assert.Contains(t, response.Body.String(), "{\"id\":\"10\",\"content\":\"Test Message 2\",\"author\":\"test author 2\",\"createdAt\":\"2017-08-15T00:00:00Z\",\"palindrome\":false}")
				assert.Equal(t, http.StatusOK, response.Code)
			},
		},
	}

	request, _ := http.NewRequest(http.MethodGet, "/messages", nil)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			//get in memory "database"
			messageRepository, _ := persistence.NewMemoryRepository()
			//database fixture
			testCase.preload(messageRepository)
			//setup a new message controller with the in memory database
			messageController := NewMessageController(messageRepository)
			//handler function to test
			handler := messageController.ListMessages

			response := httptest.NewRecorder()
			handler(response, request)
			testCase.checker(t, response)
		})
	}
}

//------------------------------- Create -----------------------------------------
func Test_Create(t *testing.T) {
	testCases := []struct {
		name    string
		preload func(memoryRepository *persistence.MemoryRepository)
		request *http.Request
		checker func(t *testing.T, response *httptest.ResponseRecorder)
	}{
		{
			name:    "Success path - content is not palindromic and empty repository",
			preload: func(memoryRepository *persistence.MemoryRepository) {},
			request: func() *http.Request {
				body := fmt.Sprint(`
                    {
                        "content": "Not a palindrome",
                        "author": "Author 1",
                        "createdAt": "2019-05-20T12:23:36.138Z"
                    }`,
				)
				request, _ := http.NewRequest(http.MethodPost, "/messages", strings.NewReader(body))
				return request
			}(),
			checker: func(t *testing.T, response *httptest.ResponseRecorder) {
				assert.Equal(t, "{\"id\":\"1\",\"content\":\"Not a palindrome\",\"author\":\"Author 1\",\"createdAt\":\"2019-05-20T12:23:36.138Z\",\"palindrome\":false}\n",
					response.Body.String())
				assert.Equal(t, http.StatusOK, response.Code)
			},
		},
		{
			name:    "Success path - content is palindromic and empty repository",
			preload: func(memoryRepository *persistence.MemoryRepository) {},
			request: func() *http.Request {
				body := fmt.Sprint(`
                    {
                        "content": "pal ind rome 12 3 21! emordnilap",
                        "author": "Author 1",
                        "createdAt": "2019-05-20T12:23:36.138Z"
                    }`,
				)
				request, _ := http.NewRequest(http.MethodPost, "/messages", strings.NewReader(body))
				return request
			}(),
			checker: func(t *testing.T, response *httptest.ResponseRecorder) {
				assert.Equal(t, "{\"id\":\"1\",\"content\":\"pal ind rome 12 3 21! emordnilap\",\"author\":\"Author 1\",\"createdAt\":\"2019-05-20T12:23:36.138Z\",\"palindrome\":true}\n",
					response.Body.String())
				assert.Equal(t, http.StatusOK, response.Code)
			},
		},
		{
			name:    "Success path - content is palindromic, contains content only and empty repository",
			preload: func(memoryRepository *persistence.MemoryRepository) {},
			request: func() *http.Request {
				body := fmt.Sprint(`
                    {
                        "content": "pal ind rome 12 3 21! emordnilap"
                    }`,
				)
				request, _ := http.NewRequest(http.MethodPost, "/messages", strings.NewReader(body))
				return request
			}(),
			checker: func(t *testing.T, response *httptest.ResponseRecorder) {
				assert.Equal(t, "{\"id\":\"1\",\"content\":\"pal ind rome 12 3 21! emordnilap\",\"palindrome\":true}\n",
					response.Body.String())
				assert.Equal(t, http.StatusOK, response.Code)
			},
		},
		{
			name:    "Fail path - empty body",
			preload: func(memoryRepository *persistence.MemoryRepository) {},
			request: func() *http.Request {
				request, _ := http.NewRequest(http.MethodPost, "/messages", strings.NewReader(""))
				return request
			}(),
			checker: func(t *testing.T, response *httptest.ResponseRecorder) {
				assert.Equal(t, "{\"message\":\"Could not decode request body: EOF\"}\n",
					response.Body.String())
				assert.Equal(t, http.StatusBadRequest, response.Code)
			},
		},
		{
			name:    "Fail path - empty JSON",
			preload: func(memoryRepository *persistence.MemoryRepository) {},
			request: func() *http.Request {
				request, _ := http.NewRequest(http.MethodPost, "/messages", strings.NewReader("{}"))
				return request
			}(),
			checker: func(t *testing.T, response *httptest.ResponseRecorder) {
				assert.Equal(t, "{\"message\":[\"Content must be between 1 and 256 characters long. Got NULL instead\"]}\n",
					response.Body.String())
				assert.Equal(t, http.StatusBadRequest, response.Code)
			},
		},
		{
			name:    "Fail path - message creation time can't be parsed",
			preload: func(memoryRepository *persistence.MemoryRepository) {},
			request: func() *http.Request {
				body := fmt.Sprint(`
                    {
                        "author": "Author 1",
                        "createdAt": "This is wrong!"
                    }`,
				)
				request, _ := http.NewRequest(http.MethodPost, "/messages", strings.NewReader(body))
				return request
			}(),
			checker: func(t *testing.T, response *httptest.ResponseRecorder) {
				assert.Contains(t, response.Body.String(), "Could not decode request body: parsing time")
				assert.Equal(t, http.StatusBadRequest, response.Code)
			},
		},
		{
			name:    "Fail path - message content missing",
			preload: func(memoryRepository *persistence.MemoryRepository) {},
			request: func() *http.Request {
				body := fmt.Sprint(`
                    {
                        "author": "Author 1",
                        "createdAt": "2019-05-20T12:23:36.138Z"
                    }`,
				)
				request, _ := http.NewRequest(http.MethodPost, "/messages", strings.NewReader(body))
				return request
			}(),
			checker: func(t *testing.T, response *httptest.ResponseRecorder) {
				assert.Equal(t, "{\"message\":[\"Content must be between 1 and 256 characters long. Got NULL instead\"]}\n",
					response.Body.String())
				assert.Equal(t, http.StatusBadRequest, response.Code)
			},
		},
		{
			name:    "Fail path - message content too short",
			preload: func(memoryRepository *persistence.MemoryRepository) {},
			request: func() *http.Request {
				body := fmt.Sprint(`
                    {
                        "content": "",
                        "author": "Author 1",
                        "createdAt": "2019-05-20T12:23:36.138Z"
                    }`,
				)
				request, _ := http.NewRequest(http.MethodPost, "/messages", strings.NewReader(body))
				return request
			}(),
			checker: func(t *testing.T, response *httptest.ResponseRecorder) {
				assert.Equal(t, "{\"message\":[\"Content must be between 1 and 256 characters long. Got 0 instead\"]}\n",
					response.Body.String())
				assert.Equal(t, http.StatusBadRequest, response.Code)
			},
		},
		{
			name:    "Fail path - message content too long",
			preload: func(memoryRepository *persistence.MemoryRepository) {},
			request: func() *http.Request {
				body := fmt.Sprint(`
                    {
                        "content": "01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789",
                        "author": "Author 1",
                        "createdAt": "2019-05-20T12:23:36.138Z"
                    }`,
				)
				request, _ := http.NewRequest(http.MethodPost, "/messages", strings.NewReader(body))
				return request
			}(),
			checker: func(t *testing.T, response *httptest.ResponseRecorder) {
				assert.Equal(t, "{\"message\":[\"Content must be between 1 and 256 characters long. Got 260 instead\"]}\n",
					response.Body.String())
				assert.Equal(t, http.StatusBadRequest, response.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			//get in memory "database"
			messageRepository, _ := persistence.NewMemoryRepository()
			//database fixture
			testCase.preload(messageRepository)
			//setup a new message controller with the in memory database
			messageController := NewMessageController(messageRepository)
			//handler function to test
			handler := messageController.CreateMessage

			response := httptest.NewRecorder()
			handler(response, testCase.request)
			testCase.checker(t, response)
		})
	}
}

//------------------------------- Get --------------------------------------------
//TODO
//------------------------------- Update -----------------------------------------
//TODO
//------------------------------- Delete -----------------------------------------
//TODO
//------------------------------- Validation -------------------------------------
//TODO

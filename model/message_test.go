package model

import (
	"encoding/json"
	"testing"
	"time"
)

func getNewString(str string) *string {
	return &str
}

func getNewMessageTime(t time.Time) *MessageTime {
	mt := MessageTime(t)
	return &mt
}

func marshal(inputStruck interface{}) string {
	errsStr, _ := json.Marshal(inputStruck)
	return string(errsStr)
}

func TestMessageValidation(t *testing.T) {
	testCases := []struct {
		messageRequest MessageRequest
		expectedLength int
	}{
		{
			MessageRequest{},
			1,
		},
		{
			MessageRequest{
				Author:    getNewString("test author"),
				CreatedAt: &MessageTime{},
			},
			1,
		},
		{
			MessageRequest{
				Author:    getNewString("test author"),
				CreatedAt: getNewMessageTime(time.Date(2016, time.August, 15, 0, 0, 0, 0, time.UTC)),
			},
			1,
		},
		{
			MessageRequest{
				Content: getNewString(""),
			},
			1,
		},
		{
			MessageRequest{
				Content:   getNewString(""),
				Author:    getNewString("test author"),
				CreatedAt: &MessageTime{},
			},
			1,
		},
		{
			MessageRequest{
				Content:   getNewString("0"),
				Author:    getNewString("test author"),
				CreatedAt: getNewMessageTime(time.Date(2016, time.August, 15, 0, 0, 0, 0, time.UTC)),
			},
			0,
		},
		{
			MessageRequest{
				Content:   getNewString("012345"),
			},
			0,
		},
		{
			MessageRequest{
				Content:   getNewString("01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789"),
				Author:    getNewString("test author"),
				CreatedAt: getNewMessageTime(time.Date(2016, time.August, 15, 0, 0, 0, 0, time.UTC)),
			},
			1,
		},
	}

	for _, testCase := range testCases {
		errs := testCase.messageRequest.Validate()
		if len(errs.Messages) != testCase.expectedLength {
			t.Errorf("For %s expected errors length of %d but got %d with %s",
				marshal(testCase.messageRequest),
				testCase.expectedLength,
				len(errs.Messages),
				marshal(errs),
			)
		}
	}
}

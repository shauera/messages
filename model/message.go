package model

import (
//"time"
)

// MessageRequest is a word, sentence or phrase written by an author
// on a specific date and timeproduct in the store.
// It is used to describe the animals available in the store.
//
// swagger:model
type MessageRequest struct {
	// The contet of the message.
	//
	// required: true
	// example: To be, or not to be: that is the question
	Content *string `json:"content,omitempty" bson:"content,omitempty"`

	// The author of the message.
	//
	// required: false
	// example: William Shakespeare
	Author *string `json:"author,omitempty" bson:"author,omitempty"`

	// The date and time when the message was created.
	//
	// required: false
	// example: 1599-01-03T07:30:30.457Z
	CreatedAt *MessageTime `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}

// MessageResponse is a word, sentence or phrase written by an author
// on a specific date and timeproduct in the store.
// It is used to describe the animals available in the store.
//
// swagger:model
type MessageResponse struct {
	// The id of the message - can't be explicitly set.
	ID interface{} `json:"id" bson:"_id"`

	// The contet of the message.
	Content *string `json:"content,omitempty" bson:"content,omitempty"`

	// The author of the message.
	Author *string `json:"author,omitempty" bson:"author,omitempty"`

	// The date and time when the message was created.
	CreatedAt *MessageTime `json:"createdAt,omitempty" bson:"createdAt,omitempty"`

	// Indicates if the message content is a palindrome.
	// This is a calculated field that can't be explicitly set.
	Palindrome bool `json:"palindrome" bson:"palindrome"`
}

// MessageResponses - a collection of MessageResponse objects
//
// swagger:model
type MessageResponses []MessageResponse

// TODO - add validation code

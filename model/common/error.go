package common

type ErrorResponse struct {
	Message   string `json:"message" bson:"message"`
}
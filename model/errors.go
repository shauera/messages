package model

//ErrorResponse - template for rendering errors in HTTP responses
type ErrorResponse struct {
	Message string `json:"message"`
}

//ValidationErrorsResponse - template for rendering errors in HTTP responses all validation errors for a specific request 
type ValidationErrorsResponse struct {
	Message []string `json:"message"`
}
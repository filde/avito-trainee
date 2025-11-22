package models

type ErrorType struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error *ErrorType `json:"error"`
}

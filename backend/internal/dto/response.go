package dto

type APIResponse struct {
	Data  interface{} `json:"data"`
	Error *APIError   `json:"error"`
	Meta  *Meta       `json:"meta,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Meta struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
	Total   int `json:"total"`
}

func Success(data interface{}) APIResponse {
	return APIResponse{Data: data, Error: nil}
}

func SuccessWithMeta(data interface{}, meta *Meta) APIResponse {
	return APIResponse{Data: data, Error: nil, Meta: meta}
}

func ErrorResponse(code, message string) APIResponse {
	return APIResponse{Data: nil, Error: &APIError{Code: code, Message: message}}
}

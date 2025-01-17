package models

type ErrorResponse struct {
	Err string `json:"err"`
}

func MakeErrorResponse(err error) ErrorResponse {
	return ErrorResponse{Err: err.Error()}
}

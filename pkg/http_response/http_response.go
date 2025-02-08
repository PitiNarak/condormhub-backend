package http_response

type HttpResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func SuccessResponse(message string, data any) *HttpResponse {
	return &HttpResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

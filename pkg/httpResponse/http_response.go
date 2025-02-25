package httpResponse

type HttpResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	Pagination any    `json:"pagination"`
	Data       any    `json:"data"`
}

func SuccessResponse(message string, data any) *HttpResponse {
	return &HttpResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func SuccessPageResponse(message string, data, pagination any) *HttpResponse {
	return &HttpResponse{
		Success:    true,
		Message:    message,
		Pagination: pagination,
		Data:       data,
	}
}

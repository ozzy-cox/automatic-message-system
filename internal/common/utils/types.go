package utils

type JsonResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Error   *string     `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

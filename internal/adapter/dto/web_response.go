package dto

type WebResponse[T any] struct {
	Total   *int64 `json:"total,omitempty"`
	Data    T      `json:"data,omitempty"`
	Message string `json:"message"`
}

func NewWebResponse[T any](data T, message string, total *int64) *WebResponse[T] {
	return &WebResponse[T]{
		Total:   total,
		Data:    data,
		Message: message,
	}
}

package response

type APIMessage struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

type ValidationError struct {
	Code    int               `json:"code"`
	Type    string            `json:"type"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

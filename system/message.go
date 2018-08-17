package system

type FormatMessage struct {
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

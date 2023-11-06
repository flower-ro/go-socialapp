package ws

type Message struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Result  any    `json:"result"`
}

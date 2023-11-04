package whatsapp

var Broadcast = make(chan BroadcastMessage)

type BroadcastMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Result  any    `json:"result"`
}

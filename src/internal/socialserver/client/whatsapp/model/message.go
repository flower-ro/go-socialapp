package model

// Reaction
type ReactionRequest struct {
	MessageID string `json:"message_id" form:"message_id"`
	Phone     string `json:"phone" form:"phone"`
	Emoji     string `json:"emoji" form:"emoji"`
}
type ReactionResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
}

// Revoke
type RevokeRequest struct {
	MessageID string `json:"message_id" uri:"message_id"`
	Phone     string `json:"phone" form:"phone"`
}
type RevokeResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
}

// Update
type UpdateMessageRequest struct {
	MessageID string `json:"message_id" uri:"message_id"`
	Message   string `json:"message" form:"message"`
	Phone     string `json:"phone" form:"phone"`
}
type UpdateMessageResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
}

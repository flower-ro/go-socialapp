package service

import "context"
import "go-socialapp/internal/socialserver/client/whatsapp/model"

type IMessageService interface {
	ReactMessage(ctx context.Context, request model.ReactionRequest) (response model.ReactionResponse, err error)
	RevokeMessage(ctx context.Context, request model.RevokeRequest) (response model.RevokeResponse, err error)
	UpdateMessage(ctx context.Context, request model.UpdateMessageRequest) (response model.UpdateMessageResponse, err error)
}

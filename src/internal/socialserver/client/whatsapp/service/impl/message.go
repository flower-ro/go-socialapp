package services

import (
	"context"
	"fmt"
	"go-socialapp/internal/pkg/third-party/whatsapp"
	"go-socialapp/internal/socialserver/client/whatsapp/model"
	"go-socialapp/internal/socialserver/client/whatsapp/service/impl/validations"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
	"time"
)

type serviceMessage struct {
	waClient *whatsapp.WaClient
}

//var messageSrv *serviceMessage
//
//func GetMessageService(waCli *whatsmeow.Client) *serviceMessage {
//	if messageSrv != nil {
//		return messageSrv
//	}
//	messageSrv = newMessageService(waCli)
//	return messageSrv
//}

func NewMessageService(waClient *whatsapp.WaClient) *serviceMessage {
	return &serviceMessage{
		waClient: waClient,
	}
}

func (service serviceMessage) ReactMessage(ctx context.Context, request model.ReactionRequest) (response model.ReactionResponse, err error) {
	if err = validations.ValidateReactMessage(ctx, request); err != nil {
		return response, err
	}
	dataWaRecipient, err := service.waClient.ValidateJidWithLogin(request.Phone)
	if err != nil {
		return response, err
	}

	msg := &waProto.Message{
		ReactionMessage: &waProto.ReactionMessage{
			Key: &waProto.MessageKey{
				FromMe:    proto.Bool(true),
				Id:        proto.String(request.MessageID),
				RemoteJid: proto.String(dataWaRecipient.String()),
			},
			Text:              proto.String(request.Emoji),
			SenderTimestampMs: proto.Int64(time.Now().UnixMilli()),
		},
	}
	ts, err := service.waClient.WaCli.SendMessage(ctx, dataWaRecipient, msg)
	if err != nil {
		return response, err
	}

	response.MessageID = ts.ID
	response.Status = fmt.Sprintf("Reaction sent to %s (server timestamp: %s)", request.Phone, ts)
	return response, nil
}

func (service serviceMessage) RevokeMessage(ctx context.Context, request model.RevokeRequest) (response model.RevokeResponse, err error) {
	if err = validations.ValidateRevokeMessage(ctx, request); err != nil {
		return response, err
	}
	dataWaRecipient, err := service.waClient.ValidateJidWithLogin(request.Phone)
	if err != nil {
		return response, err
	}

	ts, err := service.waClient.WaCli.SendMessage(context.Background(), dataWaRecipient, service.waClient.WaCli.BuildRevoke(dataWaRecipient, types.EmptyJID, request.MessageID))
	if err != nil {
		return response, err
	}

	response.MessageID = ts.ID
	response.Status = fmt.Sprintf("Revoke success %s (server timestamp: %s)", request.Phone, ts)
	return response, nil
}

func (service serviceMessage) UpdateMessage(ctx context.Context, request model.UpdateMessageRequest) (response model.UpdateMessageResponse, err error) {
	if err = validations.ValidateUpdateMessage(ctx, request); err != nil {
		return response, err
	}

	dataWaRecipient, err := service.waClient.ValidateJidWithLogin(request.Phone)
	if err != nil {
		return response, err
	}

	msg := &waProto.Message{Conversation: proto.String(request.Message)}
	ts, err := service.waClient.WaCli.SendMessage(context.Background(), dataWaRecipient, service.waClient.WaCli.BuildEdit(dataWaRecipient, request.MessageID, msg))
	if err != nil {
		return response, err
	}

	response.MessageID = ts.ID
	response.Status = fmt.Sprintf("Update message success %s (server timestamp: %s)", request.Phone, ts)
	return response, nil
}

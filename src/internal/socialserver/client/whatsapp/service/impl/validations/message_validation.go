package validations

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/marmotedu/errors"
	"go-socialapp/internal/pkg/code"
	"go-socialapp/internal/socialserver/client/whatsapp/model"
)

func ValidateRevokeMessage(ctx context.Context, request model.RevokeRequest) error {
	err := validation.ValidateStructWithContext(ctx, &request,
		validation.Field(&request.Phone, validation.Required),
		validation.Field(&request.MessageID, validation.Required),
	)

	if err != nil {
		return errors.WithCode(code.ValidationError, err.Error())
	}

	return nil
}

func ValidateUpdateMessage(ctx context.Context, request model.UpdateMessageRequest) error {
	err := validation.ValidateStructWithContext(ctx, &request,
		validation.Field(&request.Phone, validation.Required),
		validation.Field(&request.MessageID, validation.Required),
		validation.Field(&request.Message, validation.Required),
	)

	if err != nil {
		return errors.WithCode(code.ValidationError, err.Error())
	}

	return nil
}

func ValidateReactMessage(ctx context.Context, request model.ReactionRequest) error {
	err := validation.ValidateStructWithContext(ctx, &request,
		validation.Field(&request.Phone, validation.Required),
		validation.Field(&request.MessageID, validation.Required),
		validation.Field(&request.Emoji, validation.Required),
	)

	if err != nil {
		return errors.WithCode(code.ValidationError, err.Error())
	}

	return nil
}

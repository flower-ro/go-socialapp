package validations

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/marmotedu/errors"
	"go-socialapp/internal/pkg/code"
	"go-socialapp/internal/socialserver/client/whatsapp/model"
)

func ValidateUserInfo(ctx context.Context, request model.InfoRequest) error {
	err := validation.ValidateStructWithContext(ctx, &request,
		validation.Field(&request.Phone, validation.Required),
	)

	if err != nil {
		return errors.WithCode(code.ValidationError, err.Error())
	}

	return nil
}
func ValidateUserAvatar(ctx context.Context, request model.AvatarRequest) error {
	err := validation.ValidateStructWithContext(ctx, &request,
		validation.Field(&request.Phone, validation.Required),
		validation.Field(&request.IsCommunity, validation.When(request.IsCommunity, validation.Required, validation.In(true, false))),
		validation.Field(&request.IsPreview, validation.When(request.IsPreview, validation.Required, validation.In(true, false))),
	)

	if err != nil {
		return errors.WithCode(code.ValidationError, err.Error())
	}

	return nil
}

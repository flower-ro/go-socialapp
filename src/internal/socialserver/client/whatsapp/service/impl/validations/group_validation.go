package validations

import (
	"context"
	"github.com/marmotedu/errors"
	"go-socialapp/internal/pkg/code"
	"go-socialapp/internal/socialserver/client/whatsapp/model"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ValidateJoinGroupWithLink(ctx context.Context, request model.JoinGroupWithLinkRequest) error {
	err := validation.ValidateStructWithContext(ctx, &request,
		validation.Field(&request.Link, validation.Required),
	)

	if err != nil {
		return errors.WithCode(code.ValidationError, err.Error())
	}

	return nil
}

func ValidateLeaveGroup(ctx context.Context, request model.LeaveGroupRequest) error {
	err := validation.ValidateStructWithContext(ctx, &request,
		validation.Field(&request.GroupID, validation.Required),
	)

	if err != nil {
		return errors.WithCode(code.ValidationError, err.Error())
	}

	return nil
}

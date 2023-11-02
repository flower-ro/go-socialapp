package validations

import (
	"context"
	"github.com/marmotedu/errors"
	"github.com/stretchr/testify/assert"
	"go-socialapp/internal/pkg/code"
	"go-socialapp/internal/socialserver/client/whatsapp/model"
	"testing"
)

func TestValidateUserAvatar(t *testing.T) {
	type args struct {
		request model.AvatarRequest
	}
	tests := []struct {
		name string
		args args
		err  any
	}{
		{
			name: "should success",
			args: args{request: model.AvatarRequest{
				Phone:       "1728937129312@s.whatsapp.net",
				IsPreview:   false,
				IsCommunity: false,
			}},
			err: nil,
		},
		{
			name: "should error with empty phone",
			args: args{request: model.AvatarRequest{
				Phone:       "",
				IsPreview:   false,
				IsCommunity: false,
			}},
			err: errors.WithCode(code.ValidationError, "phone: cannot be blank."),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUserAvatar(context.Background(), tt.args.request)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestValidateUserInfo(t *testing.T) {
	type args struct {
		request model.InfoRequest
	}
	tests := []struct {
		name string
		args args
		err  any
	}{
		{
			name: "should success",
			args: args{request: model.InfoRequest{
				Phone: "1728937129312@s.whatsapp.net",
			}},
			err: nil,
		},
		{
			name: "should error with empty phone",
			args: args{request: model.InfoRequest{
				Phone: "",
			}},
			err: errors.WithCode(code.ValidationError, "phone: cannot be blank."),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUserInfo(context.Background(), tt.args.request)
			assert.Equal(t, tt.err, err)
		})
	}
}

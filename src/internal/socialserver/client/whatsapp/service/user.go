package service

import (
	"context"
	"go-socialapp/internal/socialserver/client/whatsapp/model"
)

type IUserService interface {
	Info(ctx context.Context, phone string) (response model.InfoResponse, err error)
	Avatar(ctx context.Context, request model.AvatarRequest) (response model.AvatarResponse, err error)
	MyListGroups(ctx context.Context) (response model.MyListGroupsResponse, err error)
	MyPrivacySetting(ctx context.Context) (response model.MyPrivacySettingResponse, err error)
}

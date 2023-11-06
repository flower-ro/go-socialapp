package service

import (
	"context"
	"go-socialapp/internal/socialserver/client/whatsapp/model"
)

type IAppService interface {
	//Login(ctx context.Context) (response model.LoginResponse, err error)
	//GetQrCode(ctx context.Context) (<-chan whatsmeow.QRChannelItem, error)
	Logout(ctx context.Context) (err error)
	Reconnect(ctx context.Context) (err error)
	FirstDevice(ctx context.Context) (response model.DevicesResponse, err error)
	FetchDevices(ctx context.Context) (response []model.DevicesResponse, err error)
}

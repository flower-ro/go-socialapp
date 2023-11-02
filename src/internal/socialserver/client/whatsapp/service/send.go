package service

import (
	"context"
	"go-socialapp/internal/socialserver/client/whatsapp/model"
)

type ISendService interface {
	SendText(ctx context.Context, request model.MessageRequest) (response model.MessageResponse, err error)
	SendImage(ctx context.Context, request model.ImageRequest) (response model.ImageResponse, err error)
	SendFile(ctx context.Context, request model.FileRequest) (response model.FileResponse, err error)
	SendVideo(ctx context.Context, request model.VideoRequest) (response model.VideoResponse, err error)
	SendContact(ctx context.Context, request model.ContactRequest) (response model.ContactResponse, err error)
	SendLink(ctx context.Context, request model.LinkRequest) (response model.LinkResponse, err error)
	SendLocation(ctx context.Context, request model.LocationRequest) (response model.LocationResponse, err error)
}

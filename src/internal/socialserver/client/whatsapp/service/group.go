package service

import (
	"context"
	"go.mau.fi/whatsmeow/types"
)
import "go-socialapp/internal/socialserver/client/whatsapp/model"

type IGroupService interface {
	JoinGroupWithLink(ctx context.Context, request model.JoinGroupWithLinkRequest) (groupID string, err error)
	LeaveGroup(ctx context.Context, request model.LeaveGroupRequest) (err error)
	CreateGroup(name string, participants []types.JID) error
}

package service

import "context"
import "go-socialapp/internal/socialserver/client/whatsapp/model"

type IGroupService interface {
	JoinGroupWithLink(ctx context.Context, request model.JoinGroupWithLinkRequest) (groupID string, err error)
	LeaveGroup(ctx context.Context, request model.LeaveGroupRequest) (err error)
}

package services

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"go-socialapp/internal/pkg/third-party/whatsapp"
	"go-socialapp/internal/socialserver/client/whatsapp/model"
	"go-socialapp/internal/socialserver/client/whatsapp/service/impl/validations"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

type groupService struct {
	waCli *whatsmeow.Client
}

//var groupSrv *groupService
//
////func GetGroupService(waCli *whatsmeow.Client) *groupService {
////	if groupSrv != nil {
////		return groupSrv
////	}
////	groupSrv = newGroupService(waCli)
////	return groupSrv
////}

func NewGroupService(waCli *whatsmeow.Client) *groupService {
	return &groupService{
		waCli: waCli,
	}
}

func (service groupService) JoinGroupWithLink(ctx context.Context, request model.JoinGroupWithLinkRequest) (groupID string, err error) {
	if err = validations.ValidateJoinGroupWithLink(ctx, request); err != nil {
		return groupID, err
	}
	err = whatsapp.MustLogin(service.waCli)
	if err != nil {
		return "", err
	}

	jid, err := service.waCli.JoinGroupWithLink(request.Link)
	if err != nil {
		return
	}
	return jid.String(), nil
}

func (service groupService) LeaveGroup(ctx context.Context, request model.LeaveGroupRequest) (err error) {
	if err = validations.ValidateLeaveGroup(ctx, request); err != nil {
		return err
	}

	JID, err := whatsapp.ValidateJidWithLogin(service.waCli, request.GroupID)
	if err != nil {
		return err
	}

	return service.waCli.LeaveGroup(JID)
}

func (service groupService) CreateGroup(name string, members []string) error {
	ps := make([]types.JID, 0, len(members))
	for _, member := range members {
		p := types.NewJID(member, "")
		ps = append(ps, p)
	}
	req := whatsmeow.ReqCreateGroup{
		Name:         name,
		Participants: ps,
	}
	group, err := service.waCli.CreateGroup(req)
	spew.Dump(group)
	if err != nil {
		return err
	}

	return nil

}

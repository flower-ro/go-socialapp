package services

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"github.com/marmotedu/errors"
	"go-socialapp/internal/pkg/third-party/whatsapp"
	utils "go-socialapp/internal/pkg/util"
	"go-socialapp/internal/socialserver/client/whatsapp/model"
	"go-socialapp/internal/socialserver/client/whatsapp/service/impl/validations"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

type groupService struct {
	waClient *whatsapp.WaClient
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

func NewGroupService(waClient *whatsapp.WaClient) *groupService {
	return &groupService{
		waClient: waClient,
	}
}

func (service groupService) JoinGroupWithLink(ctx context.Context, request model.JoinGroupWithLinkRequest) (groupID string, err error) {
	if err = validations.ValidateJoinGroupWithLink(ctx, request); err != nil {
		return groupID, err
	}
	err = service.waClient.MustLogin()
	if err != nil {
		return "", err
	}

	jid, err := service.waClient.WaCli.JoinGroupWithLink(request.Link)
	if err != nil {
		return
	}
	return jid.String(), nil
}

func (service groupService) LeaveGroup(ctx context.Context, request model.LeaveGroupRequest) (err error) {
	if err = validations.ValidateLeaveGroup(ctx, request); err != nil {
		return err
	}

	JID, err := service.waClient.ValidateJidWithLogin(request.GroupID)
	if err != nil {
		return err
	}

	return service.waClient.WaCli.LeaveGroup(JID)
}

func (service groupService) CreateGroup(name string, participants []types.JID) error {

	device, err := service.waClient.Db.GetFirstDevice()
	if err != nil {
		return errors.Wrap(err, "")
	}
	for _, j := range participants {
		tmp := j
		err = device.Contacts.PutContactName(tmp, "", utils.GenerateRandomString(4))
		if err != nil {
			return errors.Wrap(err, "")
		}
	}
	devices, err := device.Contacts.GetAllContacts()
	spew.Dump("---------devices=", devices)
	spew.Dump("---------err=", err)

	req := whatsmeow.ReqCreateGroup{
		Name:         name,
		Participants: participants,
	}
	group, err := service.waClient.WaCli.CreateGroup(req)
	if err != nil {
		return errors.Wrap(err, "")
	}
	param := map[types.JID]whatsmeow.ParticipantChange{}
	for _, p := range participants {
		tmp := p
		param[tmp] = whatsmeow.ParticipantChangeAdd

	}
	node, err := service.waClient.WaCli.UpdateGroupParticipants(group.JID, param)
	spew.Dump(node)
	if err != nil {
		return errors.Wrap(err, "")
	}

	return nil

}

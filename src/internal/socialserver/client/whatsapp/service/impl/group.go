package services

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"github.com/marmotedu/errors"
	"go-socialapp/internal/pkg/third-party/whatsapp"
	"go-socialapp/internal/socialserver/client/whatsapp/model"
	"go-socialapp/internal/socialserver/client/whatsapp/service/impl/validations"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
)

type groupService struct {
	waCli *whatsmeow.Client
	db    *sqlstore.Container
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

func NewGroupService(waCli *whatsmeow.Client, db *sqlstore.Container) *groupService {
	return &groupService{
		waCli: waCli,
		db:    db,
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

func (service groupService) CreateGroup(name string, participants []types.JID) error {

	//device, err := service.db.GetFirstDevice()
	//if err != nil {
	//	return errors.Wrap(err, "")
	//}
	//for _, j := range participants {
	//	tmp := j
	//	err = device.Contacts.PutContactName(tmp, "", utils.GenerateRandomString(4))
	//	if err != nil {
	//		return errors.Wrap(err, "")
	//	}
	//}
	//devices, err := device.Contacts.GetAllContacts()
	//spew.Dump("---------devices=", devices)
	//spew.Dump("---------err=", err)

	req := whatsmeow.ReqCreateGroup{
		Name:         name,
		Participants: participants,
	}
	group, err := service.waCli.CreateGroup(req)
	if err != nil {
		return errors.Wrap(err, "")
	}
	param := map[types.JID]whatsmeow.ParticipantChange{}
	for _, p := range participants {
		param[p] = whatsmeow.ParticipantChangeAdd

	}
	node, err := service.waCli.UpdateGroupParticipants(group.JID, param)
	spew.Dump(node)
	if err != nil {
		return errors.Wrap(err, "")
	}

	return nil

}

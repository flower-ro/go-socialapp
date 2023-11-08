package group

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"github.com/marmotedu/errors"
	"github.com/marmotedu/iam/pkg/log"
	whatsappBase "go-socialapp/internal/pkg/third-party/whatsapp"
	"go-socialapp/internal/socialserver/cache/loggedin"
	"go-socialapp/internal/socialserver/model/v1/request"
	"go-socialapp/internal/socialserver/store"
	"go-socialapp/internal/socialserver/ws"
	transcationalDB "go-socialapp/pkg/db"
	"go.mau.fi/whatsmeow/types"
	"strings"
)

type GroupSrv interface {
	Create(ctx context.Context, req request.GroupCreateReq) error
}

type groupService struct {
	transcationalDB.TxGenerate
	store store.Factory
}

var _ GroupSrv = (*groupService)(nil)

var groupSrv *groupService

func GetGroup(store store.Factory) *groupService {
	if groupSrv != nil {
		return groupSrv
	}
	groupSrv = &groupService{
		TxGenerate: store.GetTxGenerate(),
		store:      store}
	return groupSrv
}

func (a *groupService) Create(ctx context.Context, req request.GroupCreateReq) error {

	waApi, err := loggedin.WaApiCache.Get(req.Creator)
	if err != nil {
		return err
	}
	if len(req.Member) <= 0 {
		return errors.New("members cannot be 0 when create group ")
	}
	var last = make([]string, 0, len(req.Member))
	for _, member := range req.Member {
		if !strings.HasPrefix(member, "86") {
			correct := "86" + member
			last = append(last, correct)
		}
	}
	spew.Dump("------last,", last)
	result, err := waApi.General().IsOnWhatsApp(last)
	if err != nil {
		return errors.Wrap(err, "")
	}
	spew.Dump("------IsOnWhatsApp,", result)
	if len(result) > 0 {
		return nil
	}
	members := make([]types.JID, 0, len(result))
	for _, one := range result {
		members = append(members, one.JID)
	}
	spew.Dump("------members,", members)
	go func() {
		log.Infof("phone %s ,wait login", req.Creator)
		err = whatsappBase.WaitLogin(waApi.GetClient())
		if err != nil {
			log.Errorf("CreateGroup:%+v", err)
			ws.Manager.BroadcastMsg(ws.Message{
				Code:    "CREATE_GROUP_FAIL",
				Message: "login fail",
			})
			return
		}
		log.Infof("phone %s , login success", req.Creator)
		err = waApi.Group().CreateGroup(req.Name, members)
		if err != nil {
			log.Errorf("CreateGroup err:%+v", err)
			ws.Manager.BroadcastMsg(ws.Message{
				Code:    "CREATE_GROUP_FAIL",
				Message: "create group err",
			})
			return
		}
		ws.Manager.BroadcastMsg(ws.Message{
			Code: "CREATE_GROUP_SUCCESS",
		})
	}()
	return nil
}

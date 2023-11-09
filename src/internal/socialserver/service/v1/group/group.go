package group

import (
	"context"
	"github.com/marmotedu/errors"
	"go-socialapp/internal/socialserver/cache/loggedin"
	"go-socialapp/internal/socialserver/model/network"
	"go-socialapp/internal/socialserver/store"
	transcationalDB "go-socialapp/pkg/db"
	"go.mau.fi/whatsmeow/types"
)

type GroupSrv interface {
	Create(ctx context.Context, req network.GroupCreateReq) error
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

func (a *groupService) Create(ctx context.Context, req network.GroupCreateReq) error {

	waApi, err := loggedin.WaApiCache.Get(req.Creator)
	if err != nil {
		return errors.Wrap(err, "")
	}
	if len(req.Member) <= 0 {
		return errors.New("members cannot be 0 when create group ")
	}
	err = waApi.GetClient().WaitLogin()
	if err != nil {
		return errors.Wrap(err, "")
	}
	result, err := waApi.General().IsOnWhatsApp(req.Member)
	if err != nil {
		return errors.Wrap(err, "")
	}
	if len(result) <= 0 {
		return errors.New("has not valid member")
	}
	members := make([]types.JID, 0, len(result))
	for _, one := range result {
		members = append(members, one.JID)
	}
	err = waApi.Group().CreateGroup(req.Name, members)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

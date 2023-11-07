package group

import (
	"context"
	"go-socialapp/internal/socialserver/cache/loggedin"
	"go-socialapp/internal/socialserver/model/v1/request"
	"go-socialapp/internal/socialserver/store"
	transcationalDB "go-socialapp/pkg/db"
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

	waClient, err := loggedin.WaClientCache.Get(req.Creator)
	if err != nil {
		return err
	}

	return waClient.Group().CreateGroup(req.Name, req.Member)
}

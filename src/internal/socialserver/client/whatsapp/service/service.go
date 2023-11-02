package service

import (
	services "go-socialapp/internal/socialserver/client/whatsapp/service/impl"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

type thirdClient struct {
	waCli *whatsmeow.Client
	db    *sqlstore.Container
}

func NewThirdClient(waCli *whatsmeow.Client, db *sqlstore.Container) *thirdClient {
	return &thirdClient{
		waCli: waCli,
		db:    db,
	}
}

func (t *thirdClient) App() IAppService {
	return services.GetAppService(t.waCli, t.db)
}

func (t *thirdClient) Group() IGroupService {
	return services.GetGroupService(t.waCli)
}

func (t *thirdClient) Message() IMessageService {
	return services.GetMessageService(t.waCli)
}

func (t *thirdClient) Send() ISendService {
	return services.GetSendService(t.waCli, t.App())
}

func (t *thirdClient) User() IUserService {
	return services.GetUserService(t.waCli)
}

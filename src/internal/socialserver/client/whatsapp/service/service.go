package service

import (
	services "go-socialapp/internal/socialserver/client/whatsapp/service/impl"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"time"
)

type waApi struct {
	lastOperationTime time.Time
	waCli             *whatsmeow.Client
	db                *sqlstore.Container
}

func NewWaApi(waCli *whatsmeow.Client, db *sqlstore.Container) *waApi {
	api := &waApi{
		waCli: waCli,
		db:    db,
	}
	api.UpdateLastOperationTime()
	return api
}
func (t *waApi) GetClient() *whatsmeow.Client {
	return t.waCli
}

func (t *waApi) UpdateLastOperationTime() {
	t.lastOperationTime = time.Now()
}

func (t *waApi) App() IAppService {
	return services.NewAppService(t.waCli, t.db)
}

func (t *waApi) Group() IGroupService {
	return services.NewGroupService(t.waCli, t.db)
}

func (t *waApi) Message() IMessageService {
	return services.NewMessageService(t.waCli)
}

func (t *waApi) Send() ISendService {
	return services.NewSendService(t.waCli, t.App())
}

func (t *waApi) User() IUserService {
	return services.NewUserService(t.waCli)
}

func (t *waApi) General() IGeneralService {
	return services.NewGeneralService(t.waCli)
}

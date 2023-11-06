package service

import (
	services "go-socialapp/internal/socialserver/client/whatsapp/service/impl"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"time"
)

type waClient struct {
	lastOperationTime time.Time
	waCli             *whatsmeow.Client
	db                *sqlstore.Container
}

func NewClient(waCli *whatsmeow.Client, db *sqlstore.Container) *waClient {
	client := &waClient{
		waCli: waCli,
		db:    db,
	}
	client.UpdateLastOperationTime()
	return client
}
func (t *waClient) UpdateLastOperationTime() {
	t.lastOperationTime = time.Now()
}

func (t *waClient) App() IAppService {
	return services.NewAppService(t.waCli, t.db)
}

func (t *waClient) Group() IGroupService {
	return services.NewGroupService(t.waCli)
}

func (t *waClient) Message() IMessageService {
	return services.NewMessageService(t.waCli)
}

func (t *waClient) Send() ISendService {
	return services.NewSendService(t.waCli, t.App())
}

func (t *waClient) User() IUserService {
	return services.NewUserService(t.waCli)
}

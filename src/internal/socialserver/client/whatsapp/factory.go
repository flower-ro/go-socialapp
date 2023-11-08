package whatsapp

import (
	"go-socialapp/internal/socialserver/client/whatsapp/service"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
)

type Factory interface {
	App() service.IAppService
	Group() service.IGroupService
	Message() service.IMessageService
	Send() service.ISendService
	User() service.IUserService
	General() service.IGeneralService
	UpdateLastOperationTime()
	GetClient() *whatsmeow.Client
}

func NewFactory(waCli *whatsmeow.Client, db *sqlstore.Container) Factory {
	return service.NewWaApi(waCli, db)
}

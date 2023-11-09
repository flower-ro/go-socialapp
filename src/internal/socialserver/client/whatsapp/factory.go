package whatsapp

import (
	"go-socialapp/internal/pkg/third-party/whatsapp"
	"go-socialapp/internal/socialserver/client/whatsapp/service"
)

type Factory interface {
	App() service.IAppService
	Group() service.IGroupService
	Message() service.IMessageService
	Send() service.ISendService
	User() service.IUserService
	General() service.IGeneralService
	UpdateLastOperationTime()
	GetClient() *whatsapp.WaClient
}

func NewFactory(waClient *whatsapp.WaClient) Factory {
	return service.NewWaApi(waClient)
}

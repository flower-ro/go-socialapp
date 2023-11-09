package service

import (
	"go-socialapp/internal/pkg/third-party/whatsapp"
	services "go-socialapp/internal/socialserver/client/whatsapp/service/impl"
	"time"
)

type waApi struct {
	lastOperationTime time.Time
	waClient          *whatsapp.WaClient
}

func NewWaApi(waClient *whatsapp.WaClient) *waApi {
	api := &waApi{
		waClient: waClient,
	}
	api.UpdateLastOperationTime()
	return api
}
func (t *waApi) GetClient() *whatsapp.WaClient {
	return t.waClient
}

func (t *waApi) UpdateLastOperationTime() {
	t.lastOperationTime = time.Now()
}

func (t *waApi) App() IAppService {
	return services.NewAppService(t.waClient)
}

func (t *waApi) Group() IGroupService {
	return services.NewGroupService(t.waClient)
}

func (t *waApi) Message() IMessageService {
	return services.NewMessageService(t.waClient)
}

func (t *waApi) Send() ISendService {
	return services.NewSendService(t.waClient, t.App())
}

func (t *waApi) User() IUserService {
	return services.NewUserService(t.waClient)
}

func (t *waApi) General() IGeneralService {
	return services.NewGeneralService(t.waClient)
}

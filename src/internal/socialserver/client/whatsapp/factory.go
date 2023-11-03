package whatsapp

import "go-socialapp/internal/socialserver/client/whatsapp/service"

var client Factory

// Factory defines the tg task server platform storage interface.
type Factory interface {
	App() service.IAppService
	Group() service.IGroupService
	Message() service.IMessageService
	Send() service.ISendService
	User() service.IUserService
}

// Client return the store client instance.
func Client() Factory {
	return client
}

// SetClient set the iam store client.
func SetClient(factory Factory) {
	client = factory
}

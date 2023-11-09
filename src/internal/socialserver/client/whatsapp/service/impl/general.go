package services

import (
	"github.com/marmotedu/errors"
	"go-socialapp/internal/pkg/third-party/whatsapp"
	"go.mau.fi/whatsmeow/types"
)

type generalService struct {
	waClient *whatsapp.WaClient
}

func NewGeneralService(waClient *whatsapp.WaClient) *generalService {
	return &generalService{
		waClient: waClient,
	}
}

func (service generalService) IsOnWhatsApp(phones []string) (map[string]types.IsOnWhatsAppResponse, error) {
	var in = make(map[string]types.IsOnWhatsAppResponse, len(phones)/2)
	result, err := service.waClient.WaCli.IsOnWhatsApp(phones)
	if err != nil {
		return in, errors.Errorf("call wa IsOnWhatsApp api err %s", err.Error())
	}
	for _, one := range result {
		if one.IsIn {
			in[one.Query] = one
		}
	}
	return in, nil
}

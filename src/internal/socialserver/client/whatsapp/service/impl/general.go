package services

import (
	"github.com/marmotedu/errors"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
)

type generalService struct {
	waCli *whatsmeow.Client
}

func NewGeneralService(waCli *whatsmeow.Client) *generalService {
	return &generalService{
		waCli: waCli,
	}
}

func (service generalService) IsOnWhatsApp(phones []string) (map[string]types.IsOnWhatsAppResponse, error) {
	var in = make(map[string]types.IsOnWhatsAppResponse, len(phones)/2)
	result, err := service.waCli.IsOnWhatsApp(phones)
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

package service

import (
	"go.mau.fi/whatsmeow/types"
)

type IGeneralService interface {
	IsOnWhatsApp(phones []string) (map[string]types.IsOnWhatsAppResponse, error)
}

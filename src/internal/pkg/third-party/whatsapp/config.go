package whatsapp

import (
	"fmt"
	waProto "go.mau.fi/whatsmeow/binary/proto"
)

var (
	//PathStorages = "session"
	PathSessions    = "session/login"
	PathSessionsTmp = "session/tmp"
	PathImage       = "image"
	History         = "history"
	//DBName          = "whatsapp.db"

	AppOs            = fmt.Sprintf("AldinoKemal")
	AppVersion       = "v4.8.0"
	AppPlatform      = waProto.DeviceProps_PlatformType(1)
	WhatsappLogLevel = "ERROR"

	WhatsappAutoReplyMessage string
	WhatsappWebhook          string

	PathMedia     = "statics/media"
	PathQrCode    = "statics/qrcode"
	PathSendItems = "statics/senditems"

	WhatsappSettingMaxFileSize  int64 = 50000000  // 50MB
	WhatsappSettingMaxVideoSize int64 = 100000000 // 100MB
)

package watest

import (
	"go.mau.fi/whatsmeow/store"
	"google.golang.org/protobuf/proto"
	"strconv"

	waProto "go.mau.fi/whatsmeow/binary/proto"
)

var waVersion = store.WAVersionContainer{2, 2344, 53}

var BaseClientPayload = &waProto.ClientPayload{
	UserAgent: &waProto.ClientPayload_UserAgent{
		Platform:       waProto.ClientPayload_UserAgent_WEB.Enum(),
		ReleaseChannel: waProto.ClientPayload_UserAgent_RELEASE.Enum(),
		AppVersion:     waVersion.ProtoAppVersion(),
		Mcc:            proto.String("000"),
		Mnc:            proto.String("000"),
		OsVersion:      proto.String("0.1.0"),
		Manufacturer:   proto.String(""),
		Device:         proto.String("Desktop"),
		OsBuildNumber:  proto.String("0.1.0"),

		LocaleLanguageIso6391:       proto.String("en"),
		LocaleCountryIso31661Alpha2: proto.String("en"),
	},
	WebInfo: &waProto.ClientPayload_WebInfo{
		WebSubPlatform: waProto.ClientPayload_WebInfo_WEB_BROWSER.Enum(),
	},
	ConnectType:   waProto.ClientPayload_WIFI_UNKNOWN.Enum(),
	ConnectReason: waProto.ClientPayload_USER_ACTIVATED.Enum(),
}

//ID           *types.JID
/**
type JID struct {
	User       string
	RawAgent   uint8
	Device     uint16
	Integrator uint16
	Server     string
}
*/
func getLoginPayload(user string) *waProto.ClientPayload {
	payload := proto.Clone(BaseClientPayload).(*waProto.ClientPayload)
	payload.Username = proto.Uint64(UserInt(user))
	payload.Device = proto.Uint32(uint32(0))
	payload.Passive = proto.Bool(true)
	return payload
}

func UserInt(user string) uint64 {
	number, _ := strconv.ParseUint(user, 10, 64)
	return number
}

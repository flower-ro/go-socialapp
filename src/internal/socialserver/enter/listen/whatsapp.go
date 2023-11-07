package listen

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"github.com/marmotedu/iam/pkg/log"
	whatsappApi "go-socialapp/internal/socialserver/client/whatsapp"
	"go-socialapp/internal/socialserver/client/whatsapp/model"
	"go-socialapp/internal/socialserver/enter/ws"
	srvv1 "go-socialapp/internal/socialserver/service/v1"
	"strings"
	"time"
)
import "go-socialapp/internal/pkg/third-party/whatsapp"

var waListen *WaListen

type WaListen struct {
	listenCh chan whatsapp.BroadcastMessage
	srv      srvv1.Service
}

func InitWaListen() {
	waListen = &WaListen{
		listenCh: whatsapp.BroadcastCh,
		srv:      srvv1.GetService(),
	}
	go waListen.start()
}

func (w *WaListen) start() {

	for {
		select {
		case message, ok := <-w.listenCh:
			if !ok {
				log.Errorf("listen is exception close")
				time.Sleep(2 * time.Minute)
				continue
			}
			switch message.Code {
			case whatsapp.MessageTypeLogin:
				w.handlerLoginMessage(message)
			case whatsapp.MessageTypeLogout:
				w.handlerLogoutMessage(message)
			default:
				log.Errorf("listen wa message type is not defined")
			}
		}
	}
}

func (w *WaListen) handlerLoginMessage(message whatsapp.BroadcastMessage) error {

	factory := whatsappApi.NewFactory(message.WaClient.WaCli, message.WaClient.Db)
	var phone string
	if strings.Contains(message.Result.(string), ":") {
		strs := strings.Split(message.Result.(string), ":")
		phone = strs[0]
	} else {
		strs := strings.Split(message.Result.(string), "@")
		phone = strs[0]
	}
	spew.Dump("---message---", message.Result)
	spew.Dump("---phone---", phone)
	response, err := factory.App().FirstDevice(context.Background())

	log.Infof("FirstDevice err %v", err)
	spew.Dump("---response---", response)

	var inforesponse model.InfoResponse
	var i int
	for {
		time.Sleep(10 * time.Second)
		inforesponse, err = factory.User().Info(context.Background(), phone)
		if err == nil || i > 3 {
			break
		}
		i++
	}

	log.Infof("Info err %v", err)
	spew.Dump("---inforesponse---", inforesponse)

	ws.Manager.BroadcastMsg(ws.Message{Code: whatsapp.MessageTypeLogin, Result: message.Result})

	return nil

}

func (w *WaListen) handlerLogoutMessage(message whatsapp.BroadcastMessage) error {

	return nil
}

package listen

import (
	"github.com/marmotedu/iam/pkg/log"
	"github.com/otiai10/copy"
	"go-socialapp/internal/socialserver/cache/loggedin"
	whatsappApi "go-socialapp/internal/socialserver/client/whatsapp"
	"go-socialapp/internal/socialserver/enter/ws"
	srvv1 "go-socialapp/internal/socialserver/service/v1"
	"path/filepath"
	"strings"
	"time"
)
import "go-socialapp/internal/pkg/third-party/whatsapp"

const (
	MessageTypeLoginFail = "LOGIN_FAIL"
)

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
	var phone string
	if strings.Contains(message.Result.(string), ":") {
		strs := strings.Split(message.Result.(string), ":")
		phone = strs[0]
	} else {
		strs := strings.Split(message.Result.(string), "@")
		phone = strs[0]
	}
	err := copy.Copy(message.WaClient.Path, filepath.Join(whatsapp.PathSessions, phone+".db"))
	if err != nil {
		log.Errorf("Phone %s,copy sessionTmp %s  to session file err %s", phone, message.WaClient.Path, err.Error())
		ws.Manager.BroadcastMsg(ws.Message{Code: MessageTypeLoginFail})
		return nil
	}

	newClient, err := whatsapp.NewWaClientWithDevice(phone)
	if err != nil {
		log.Errorf("Phone %s,NewWaClientWithDevice %s err %s", phone, message.Result.(string), err.Error())
		ws.Manager.BroadcastMsg(ws.Message{Code: MessageTypeLoginFail})
		return nil
	}

	err = w.srv.Accounts().CreateOrUpdate(phone, message.Result.(string))
	if err != nil {
		log.Errorf("Phone %s,NewWaClientWithDevice err %s", phone, err.Error())
		ws.Manager.BroadcastMsg(ws.Message{Code: MessageTypeLoginFail})
		return nil
	}

	newClient.WaCli.Connect()
	factory := whatsappApi.NewFactory(newClient.WaCli, newClient.Db)
	loggedin.WaClientCache.Put(phone, factory)

	message.WaClient.WaCli.Disconnect()
	ws.Manager.BroadcastMsg(ws.Message{Code: whatsapp.MessageTypeLogin, Result: message.Result})

	return nil

}

func (w *WaListen) handlerLogoutMessage(message whatsapp.BroadcastMessage) error {

	return nil
}

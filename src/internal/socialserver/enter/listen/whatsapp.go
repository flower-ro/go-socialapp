package listen

import (
	"github.com/marmotedu/iam/pkg/log"
	srvv1 "go-socialapp/internal/socialserver/service/v1"
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
				go w.handlerLoginMessage(message)
			case whatsapp.MessageTypeLogout:
				w.handlerLogoutMessage(message)
			default:
				log.Errorf("listen wa message type is not defined")
			}
		}
	}
}

func (w *WaListen) handlerLogoutMessage(message whatsapp.BroadcastMessage) error {

	return nil
}

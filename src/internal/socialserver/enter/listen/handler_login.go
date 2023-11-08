package listen

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"github.com/marmotedu/errors"
	"github.com/marmotedu/iam/pkg/log"
	"github.com/otiai10/copy"
	"go-socialapp/internal/pkg/third-party/whatsapp"
	utils "go-socialapp/internal/pkg/util"
	"go-socialapp/internal/socialserver/cache/loggedin"
	whatsappApi "go-socialapp/internal/socialserver/client/whatsapp"
	ws2 "go-socialapp/internal/socialserver/ws"
	"path/filepath"
	"strings"
)

func (w *WaListen) handlerLoginMessage(message whatsapp.BroadcastMessage) error {
	var phone string
	var err error
	defer func() {
		if err != nil {
			message.WaClient.WaCli.Disconnect()
			log.Errorf("%+v", err)
		}
	}()
	if strings.Contains(message.Result.(string), ":") {
		strs := strings.Split(message.Result.(string), ":")
		phone = strs[0]
	} else {
		strs := strings.Split(message.Result.(string), "@")
		phone = strs[0]
	}

	err = whatsapp.WaitLogin(message.WaClient.WaCli)
	if err != nil {
		return errors.Wrap(err, " ")
	}
	newPath := filepath.Join(whatsapp.PathSessions, phone+".db")
	utils.RemoveFile(0, newPath)
	err = copy.Copy(message.WaClient.Path, newPath)
	if err != nil {
		log.Errorf("Phone %s,copy sessionTmp %s  to session file err %s", phone, message.WaClient.Path, err.Error())
		ws2.Manager.BroadcastMsg(ws2.Message{Code: MessageTypeLoginFail})
		return err
	}

	newDb, err := whatsapp.NewWaDB(newPath)
	if err != nil {
		log.Errorf("Phone %s,NewWaDB %s err %s", phone, message.Result.(string), err.Error())
		ws2.Manager.BroadcastMsg(ws2.Message{Code: MessageTypeLoginFail})
		return err
	}

	err = w.srv.Accounts().CreateOrUpdate(phone, message.Result.(string))
	if err != nil {
		log.Errorf("Phone %s,NewWaClientWithDevice err %s", phone, err.Error())
		ws2.Manager.BroadcastMsg(ws2.Message{Code: MessageTypeLoginFail})
		return err
	}
	factory := whatsappApi.NewFactory(message.WaClient.WaCli, newDb)
	loggedin.WaApiCache.Put(phone, factory)
	ws2.Manager.BroadcastMsg(ws2.Message{Code: whatsapp.MessageTypeLogin, Result: message.Result})

	spew.Dump("------------------vb--------")
	spew.Dump(factory.User().MyListGroups(context.Background()))
	return nil

}

package listen

import (
	"github.com/marmotedu/errors"
	"github.com/marmotedu/iam/pkg/log"
	"go-socialapp/internal/pkg/util/idgenerate"
	"os"
	"time"

	"github.com/otiai10/copy"
	"go-socialapp/internal/pkg/third-party/whatsapp"
	"go-socialapp/internal/socialserver/cache/loggedin"
	whatsappApi "go-socialapp/internal/socialserver/client/whatsapp"
	ws "go-socialapp/internal/socialserver/ws"
	"path/filepath"
	"strings"
)

func (w *WaListen) handlerLoginMessage(message whatsapp.BroadcastMessage) error {
	var phone, tmpFileName, newPath string
	var err error
	defer func() {
		if err != nil {
			log.Errorf("%+v", err)
			_, err = os.Stat(tmpFileName)
			if err == nil && newPath != "" {
				err = os.Rename(tmpFileName, newPath)
				if err != nil {
					log.Errorf("recover tmpFileName %s to newPath %s err %s", tmpFileName, newPath, err)
				}
			}
			message.WaClient.WaCli.Disconnect()
			ws.Manager.BroadcastMsg(ws.Message{Code: MessageTypeLoginFail})
		} else if tmpFileName != "" {
			os.Remove(tmpFileName)
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
	time.Sleep(3 * time.Minute)
	newPath = filepath.Join(whatsapp.PathSessions, phone+".db")
	tmpFileName = filepath.Join(whatsapp.PathSessions, phone+"-"+idgenerate.GetUUID36("")+".db")
	err = os.Rename(newPath, tmpFileName)
	//err = utils.RemoveFile(0, newPath)
	if err != nil {
		return errors.Wrapf(err, "Phone %s,Rename file name is %s", phone, tmpFileName)
	}
	err = copy.Copy(message.WaClient.Path, newPath)
	if err != nil {
		return errors.Wrapf(err, "Phone %s,copy sessionTmp %s  to session file ", phone, message.WaClient.Path)

	}
	//newDb, err := whatsapp.NewWaDB(newPath)
	//if err != nil {
	//	return errors.Wrapf(err, "Phone %s,NewWaDB for device %s", phone, message.Result.(string))
	//}
	newClient, err := whatsapp.NewWaClientWithDevice(phone)
	if err != nil {
		return errors.Wrapf(err, "Phone %s,NewWaClientWithDevice err ", phone)
	}
	message.WaClient.WaCli.Disconnect()
	err = whatsapp.WaitLogin(newClient.WaCli)
	if err != nil {
		return errors.Wrap(err, " new Client WaitLogin")
	}
	err = w.srv.Accounts().CreateOrUpdate(phone, message.Result.(string))
	if err != nil {
		return errors.Wrapf(err, "Phone %s,NewWaClientWithDevice", phone)
	}
	factory := whatsappApi.NewFactory(newClient.WaCli, newClient.Db)
	loggedin.WaApiCache.Put(phone, factory)
	ws.Manager.BroadcastMsg(ws.Message{Code: whatsapp.MessageTypeLogin, Result: message.Result})
	return nil

}

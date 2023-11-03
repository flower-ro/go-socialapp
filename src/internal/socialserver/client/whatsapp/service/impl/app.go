package services

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/marmotedu/errors"
	"github.com/marmotedu/iam/pkg/log"
	"time"

	"go-socialapp/internal/pkg/code"
	"go-socialapp/internal/pkg/third-party/whatsapp"
	"go-socialapp/internal/socialserver/client/whatsapp/model"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"os"
	"path/filepath"
	"strings"
)

type serviceApp struct {
	waCli *whatsmeow.Client
	db    *sqlstore.Container
}

var appservice *serviceApp

func GetAppService(waCli *whatsmeow.Client, db *sqlstore.Container) *serviceApp {
	if appservice != nil {
		return appservice
	}
	appservice = newAppService(waCli, db)
	return appservice
}

func newAppService(waCli *whatsmeow.Client, db *sqlstore.Container) *serviceApp {
	return &serviceApp{
		waCli: waCli,
		db:    db,
	}
}

func (service serviceApp) Login(_ context.Context) (response model.LoginResponse, err error) {
	if service.waCli == nil {
		return response, errors.WithCode(code.ErrWaCLI, "")
	}

	// Disconnect for reconnecting
	service.waCli.Disconnect()

	//chImage := make(chan string)
	log.Info("start get qrCode")
	ch, err := service.waCli.GetQRChannel(context.Background())
	if err != nil {
		log.Error(err.Error())
		// This error means that we're already logged in, so ignore it.
		if errors.Is(err, whatsmeow.ErrQRStoreContainsID) {
			_ = service.waCli.Connect() // just connect to websocket
			if service.waCli.IsLoggedIn() {
				return response, errors.WithCode(code.ErrAlreadyLoggedIn, "")
			}
			return response, errors.WithCode(code.ErrSessionSaved, "")
		} else {
			return response, errors.WithCode(code.ErrQrChannel, "")
		}
	} else {
		//go func() {
		for evt := range ch {
			spew.Dump("---0", evt)
			if evt.Event == "code" {
				response.Code = evt.Code
				response.Duration = evt.Timeout / time.Second / 2

				spew.Dump("---1", evt)

				//qrPath := fmt.Sprintf("%s/scan-qr-%s.png", whatsapp.PathQrCode, fiberUtils.UUIDv4())
				//err = qrcode.WriteFile(evt.Code, qrcode.Medium, 512, qrPath)
				//
				//if err != nil {
				//	log.Errorf("Error when write qr code to file: %v", err)
				//}
				//go func() {
				//	time.Sleep(response.Duration * time.Second)
				//	err := os.Remove(qrPath)
				//	if err != nil {
				//		log.Errorf("error when remove qrImage file err: %s", err.Error())
				//	}
				//}()
				//chImage <- qrPath
			} else {
				log.Errorf("error when get qrCode for event %v", evt.Event)
			}
		}
		//}()
	}

	err = service.waCli.Connect()
	if err != nil {
		return response, errors.WithCode(code.ErrReconnect, err.Error())
	}
	//response.ImagePath = <-chImage

	return response, nil
}

func (service serviceApp) Logout(_ context.Context) (err error) {
	// delete history
	files, err := filepath.Glob(fmt.Sprintf("./%s/history-*", whatsapp.PathStorages))
	if err != nil {
		return err
	}

	for _, f := range files {
		err = os.Remove(f)
		if err != nil {
			return err
		}
	}
	// delete qr images
	qrImages, err := filepath.Glob(fmt.Sprintf("./%s/scan-*", whatsapp.PathQrCode))
	if err != nil {
		return err
	}

	for _, f := range qrImages {
		err = os.Remove(f)
		if err != nil {
			return err
		}
	}

	// delete senditems
	qrItems, err := filepath.Glob(fmt.Sprintf("./%s/*", whatsapp.PathSendItems))
	if err != nil {
		return err
	}

	for _, f := range qrItems {
		if !strings.Contains(f, ".gitignore") {
			err = os.Remove(f)
			if err != nil {
				return err
			}
		}
	}

	err = service.waCli.Logout()
	return
}

func (service serviceApp) Reconnect(_ context.Context) (err error) {
	service.waCli.Disconnect()
	return service.waCli.Connect()
}

func (service serviceApp) FirstDevice(ctx context.Context) (response model.DevicesResponse, err error) {
	if service.waCli == nil {
		return response, errors.WithCode(code.ErrWaCLI, "")
	}

	devices, err := service.db.GetFirstDevice()
	if err != nil {
		return response, err
	}

	response.Device = devices.ID.String()
	if devices.PushName != "" {
		response.Name = devices.PushName
	} else {
		response.Name = devices.BusinessName
	}

	return response, nil
}

func (service serviceApp) FetchDevices(_ context.Context) (response []model.DevicesResponse, err error) {
	if service.waCli == nil {
		return response, errors.WithCode(code.ErrWaCLI, "")
	}

	devices, err := service.db.GetAllDevices()
	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		var d model.DevicesResponse
		d.Device = device.ID.String()
		if device.PushName != "" {
			d.Name = device.PushName
		} else {
			d.Name = device.BusinessName
		}

		response = append(response, d)
	}

	return response, nil
}

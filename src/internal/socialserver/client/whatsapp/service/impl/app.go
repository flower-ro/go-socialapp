package services

import (
	"context"
	"fmt"
	"github.com/marmotedu/errors"
	"go-socialapp/internal/pkg/code"
	"go-socialapp/internal/pkg/third-party/whatsapp"
	"go-socialapp/internal/socialserver/client/whatsapp/model"
	"os"
	"path/filepath"
	"strings"
)

type serviceApp struct {
	waClient *whatsapp.WaClient
}

//var appservice *serviceApp
//
//func GetAppService(waCli *whatsmeow.Client, db *sqlstore.Container) *serviceApp {
//	if appservice != nil {
//		return appservice
//	}
//	appservice = newAppService(waCli, db)
//	return appservice
//}

func NewAppService(waClient *whatsapp.WaClient) *serviceApp {
	return &serviceApp{
		waClient: waClient,
	}
}

//func (service serviceApp) GetQrCode(ctx context.Context) (<-chan whatsmeow.QRChannelItem, error) {
//	if service.waCli == nil {
//		return nil, errors.WithCode(code.ErrWaCLI, "")
//	}
//	service.waCli.Disconnect()
//
//	ch, err := service.waCli.GetQRChannel(context.Background())
//	if err != nil {
//		log.Error(err.Error())
//		// This error means that we're already logged in, so ignore it.
//		if errors.Is(err, whatsmeow.ErrQRStoreContainsID) {
//			_ = service.waCli.Connect() // just connect to websocket
//			if service.waCli.IsLoggedIn() {
//				return nil, errors.WithCode(code.ErrAlreadyLoggedIn, err.Error())
//			}
//			return nil, errors.WithCode(code.ErrSessionSaved, err.Error())
//		} else {
//			return nil, errors.WithCode(code.ErrQrChannel, err.Error())
//		}
//	}
//
//	err = service.waCli.Connect()
//	if err != nil {
//		return nil, errors.WithCode(code.ErrReconnect, err.Error())
//	}
//
//	return ch, nil
//}

//func (service serviceApp) Login(_ context.Context) (response model.LoginResponse, err error) {
//	if service.waCli == nil {
//		return response, errors.WithCode(code.ErrWaCLI, "")
//	}
//
//	// Disconnect for reconnecting
//	service.waCli.Disconnect()
//
//	//chImage := make(chan string)
//
//	ch, err := service.waCli.GetQRChannel(context.Background())
//	if err != nil {
//		log.Error(err.Error())
//		// This error means that we're already logged in, so ignore it.
//		if errors.Is(err, whatsmeow.ErrQRStoreContainsID) {
//			_ = service.waCli.Connect() // just connect to websocket
//			if service.waCli.IsLoggedIn() {
//				return response, errors.WithCode(code.ErrAlreadyLoggedIn, err.Error())
//			}
//			return response, errors.WithCode(code.ErrSessionSaved, err.Error())
//		} else {
//			return response, errors.WithCode(code.ErrQrChannel, err.Error())
//		}
//	} else {
//		go func() {
//			for evt := range ch {
//				if evt.Event == "code" {
//					response.Code = evt.Code
//					response.Duration = evt.Timeout / time.Second / 2
//				} else {
//					logrus.Error("error when get qrCode", evt.Event)
//				}
//			}
//		}()
//	}
//
//	err = service.waCli.Connect()
//	if err != nil {
//		return response, errors.WithCode(code.ErrReconnect, err.Error())
//	}
//	//response.ImagePath = <-chImage
//
//	return response, nil
//}

func (service serviceApp) Logout(_ context.Context) (err error) {
	// delete history
	files, err := filepath.Glob(fmt.Sprintf("./%s/history-*", whatsapp.PathSessions))
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

	err = service.waClient.WaCli.Logout()
	return
}

func (service serviceApp) Reconnect(_ context.Context) (err error) {
	service.waClient.WaCli.Disconnect()
	return service.waClient.WaCli.Connect()
}

func (service serviceApp) FirstDevice(ctx context.Context) (response model.DevicesResponse, err error) {
	if service.waClient.WaCli == nil {
		return response, errors.WithCode(code.ErrWaCLI, "")
	}

	devices, err := service.waClient.Db.GetFirstDevice()
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
	if service.waClient.WaCli == nil {
		return response, errors.WithCode(code.ErrWaCLI, "")
	}

	devices, err := service.waClient.Db.GetAllDevices()
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

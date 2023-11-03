package whatsapp

import (
	"fmt"
	"github.com/marmotedu/errors"
	"github.com/marmotedu/iam/pkg/log"
	_ "github.com/mattn/go-sqlite3"
	"go-socialapp/internal/pkg/code"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"time"
)

var (
	cli           *whatsmeow.Client
	historySyncID int32
	startupTime   = time.Now().Unix()
)

func InitWaDB() (*sqlstore.Container, error) {
	// Running Whatsapp
	storeContainer, err := sqlstore.New("sqlite3",
		fmt.Sprintf("file:%s/%s?_foreign_keys=off", PathStorages, DBName), nil)
	if err != nil {
		log.Errorf(err.Error())
		return nil, errors.WithCode(code.FailedConnectSqlite3, err.Error())
	}
	return storeContainer, nil
}

func InitWaCLI(storeContainer *sqlstore.Container) (*whatsmeow.Client, error) {
	device, err := storeContainer.GetFirstDevice()
	if err != nil {
		log.Errorf("Failed to get device: %v", err)
		return nil, errors.WithCode(code.FailedGetDevice, err.Error())
	}

	osName := fmt.Sprintf("%s %s", AppOs, AppVersion)
	store.DeviceProps.PlatformType = &AppPlatform
	store.DeviceProps.Os = &osName
	cli = whatsmeow.NewClient(device, nil)
	cli.EnableAutoReconnect = true
	cli.AutoTrustIdentity = true
	cli.AddEventHandler(handler)

	return cli, nil
}

package whatsapp

import (
	"fmt"
	"github.com/marmotedu/errors"
	_ "github.com/mattn/go-sqlite3"
	"go-socialapp/internal/pkg/code"
	utils "go-socialapp/internal/pkg/util"
	"go-socialapp/internal/pkg/util/idgenerate"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var lock = &sync.RWMutex{}

type WaClient struct {
	WaCli      *whatsmeow.Client
	Db         *sqlstore.Container
	FileName   string
	Path       string
	tmp        bool
	CreateTime time.Time
	user       string // ws连接获取qrcode 时候需要传送user信息，以使得，登录成功信息只会发给该user
}

func NewClientWithNoDevice() (*WaClient, error) {
	lock.Lock()
	defer lock.Unlock()
	var filePath, randomName string
	var i int
	for {
		randomName = idgenerate.GetUUID36("")
		filePath = filepath.Join(PathSessionsTmp, randomName+".db")
		_, err := os.Stat(filePath)
		if err != nil && os.IsNotExist(err) {
			break
		}

		if i > 5 {
			return nil, errors.WithCode(code.FileCreatedFail, "")
		}
		i++
	}
	client := &WaClient{
		FileName: randomName,
		Path:     filePath,
		tmp:      true,
	}
	err := client.initClient()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return client, nil
}

func NewWaClientWithDevice(waAccount string) (*WaClient, error) {
	path := filepath.Join(PathSessions, waAccount+".db")
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.WithCode(code.FileIsNotExisted, err.Error())
		}
		return nil, errors.WithCode(code.NotSureIsExisted, err.Error())
	}
	client := &WaClient{
		FileName: waAccount,
		Path:     path,
		tmp:      false,
	}
	err = client.initClient()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return client, nil
}

func (w *WaClient) initClient() error {
	err := w.initWaDB()
	if err != nil {
		return errors.Wrap(err, "")
	}
	err = w.initWaCLI()
	if err != nil {
		return errors.Wrap(err, "")
	}
	w.CreateTime = utils.GetCurrentTime()
	return nil
}

func (w *WaClient) initWaDB() error {
	// Running Whatsapp
	storeContainer, err := NewWaDB(w.Path)
	if err != nil {
		return err
	}
	w.Db = storeContainer
	return nil
}

func NewWaDB(path string) (*sqlstore.Container, error) {
	// Running Whatsapp
	storeContainer, err := sqlstore.New("sqlite3",
		fmt.Sprintf("file:%s?_foreign_keys=off", path), nil)
	if err != nil {
		return storeContainer, errors.WithCode(code.FailedConnectSqlite3, err.Error())
	}
	return storeContainer, nil
}

func (w *WaClient) initWaCLI() error {
	device, err := w.Db.GetFirstDevice()
	if err != nil {
		return errors.WithCode(code.FailedGetDevice, err.Error())
	}

	osName := fmt.Sprintf("%s %s", AppOs, AppVersion)
	store.DeviceProps.PlatformType = &AppPlatform
	store.DeviceProps.Os = &osName
	cli := whatsmeow.NewClient(device, nil)
	cli.EnableAutoReconnect = true
	cli.AutoTrustIdentity = true
	w.WaCli = cli
	cli.AddEventHandler(w.handler)
	return nil
}

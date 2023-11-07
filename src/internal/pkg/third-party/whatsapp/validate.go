package whatsapp

import (
	"fmt"
	"github.com/marmotedu/errors"
	"go-socialapp/internal/pkg/code"
	utils "go-socialapp/internal/pkg/util"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"strings"
	"time"
)

func ValidateJidWithLogin(waCli *whatsmeow.Client, jid string) (types.JID, error) {
	err := MustLogin(waCli)
	if err != nil {
		return types.JID{}, err
	}
	return ParseJID(jid)
}

func MustLogin(waCli *whatsmeow.Client) error {
	if waCli == nil {
		return errors.WithCode(code.ClientNotInitialized, "Whatsapp client is not initialized")
	}
	if !waCli.IsConnected() {
		return errors.WithCode(code.NotConnectServer, "you are not connect to services server, please reconnect")
	} else if !waCli.IsLoggedIn() {
		return errors.WithCode(code.NotLoginServer, "you are not login to services server, please login")
	}
	return nil
}

func WaitLogin(waCli *whatsmeow.Client) error {
	if waCli == nil {
		return errors.WithCode(code.ClientNotInitialized, "Whatsapp client is not initialized")
	}

	if !waCli.IsConnected() {
		waCli.Connect()
	}
	var now = utils.GetCurrentTime()
	var defaultInterval = 5 * time.Minute
	expectExpireTime := now.Add(defaultInterval)
	for {
		if now.After(expectExpireTime) || now.Equal(expectExpireTime) {
			break
		}

		if waCli.IsLoggedIn() {
			break
		}
		time.Sleep(10 * time.Second)
	}

	if !waCli.IsLoggedIn() {
		return errors.New("login fail")
	}

	return nil
}

func ParseJID(arg string) (types.JID, error) {
	if arg[0] == '+' {
		arg = arg[1:]
	}
	if !strings.ContainsRune(arg, '@') {
		return types.NewJID(arg, types.DefaultUserServer), nil
	} else {
		recipient, err := types.ParseJID(arg)
		if err != nil {
			fmt.Printf("invalid JID %s: %v", arg, err)
			return recipient, errors.WithCode(code.ErrInvalidJID, "")
		} else if recipient.User == "" {
			fmt.Printf("invalid JID %v: no server specified", arg)
			return recipient, errors.WithCode(code.ErrInvalidJID, "")
		}
		return recipient, nil
	}
}

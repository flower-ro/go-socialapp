package services

import (
	"context"
	"fmt"
	"github.com/marmotedu/errors"
	"go-socialapp/internal/pkg/code"
	"go-socialapp/internal/pkg/third-party/whatsapp"
	"go-socialapp/internal/socialserver/client/whatsapp/model"
	"go-socialapp/internal/socialserver/client/whatsapp/service/impl/validations"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"time"
)

type userService struct {
	waCli *whatsmeow.Client
}

//var userSrv *userService
//
//func GetUserService(waCli *whatsmeow.Client) *userService {
//	if userSrv != nil {
//		return userSrv
//	}
//	userSrv = newUserService(waCli)
//	return userSrv
//}

func NewUserService(waCli *whatsmeow.Client) *userService {
	return &userService{
		waCli: waCli,
	}
}

func (service userService) Info(ctx context.Context, phone string) (response model.InfoResponse, err error) {
	request := model.InfoRequest{
		Phone: phone,
	}

	err = validations.ValidateUserInfo(ctx, request)
	if err != nil {
		return response, err
	}
	var jids []types.JID
	dataWaRecipient, err := whatsapp.ValidateJidWithLogin(service.waCli, request.Phone)
	if err != nil {
		return response, err
	}

	jids = append(jids, dataWaRecipient)
	resp, err := service.waCli.GetUserInfo(jids)
	if err != nil {
		return response, err
	}

	for _, userInfo := range resp {
		var device []model.InfoResponseDataDevice
		for _, j := range userInfo.Devices {
			device = append(device, model.InfoResponseDataDevice{
				User:   j.User,
				Agent:  j.RawAgent,
				Device: whatsapp.GetPlatformName(int(j.Device)),
				Server: j.Server,
				AD:     j.ADString(),
			})
		}

		data := model.InfoResponseData{
			Status:    userInfo.Status,
			PictureID: userInfo.PictureID,
			Devices:   device,
		}
		if userInfo.VerifiedName != nil {
			data.VerifiedName = fmt.Sprintf("%v", *userInfo.VerifiedName)
		}
		response.Data = append(response.Data, data)
	}

	return response, nil
}

func (service userService) Avatar(ctx context.Context, request model.AvatarRequest) (response model.AvatarResponse, err error) {

	chanResp := make(chan model.AvatarResponse)
	chanErr := make(chan error)
	waktu := time.Now()

	go func() {
		err = validations.ValidateUserAvatar(ctx, request)
		if err != nil {
			chanErr <- err
		}
		dataWaRecipient, err := whatsapp.ValidateJidWithLogin(service.waCli, request.Phone)
		if err != nil {
			chanErr <- err
		}
		pic, err := service.waCli.GetProfilePictureInfo(dataWaRecipient, &whatsmeow.GetProfilePictureParams{
			Preview:     request.IsPreview,
			IsCommunity: request.IsCommunity,
		})
		if err != nil {
			chanErr <- err
		} else if pic == nil {
			chanErr <- errors.New("no avatar found")
		} else {
			response.URL = pic.URL
			response.ID = pic.ID
			response.Type = pic.Type

			chanResp <- response
		}
	}()

	for {
		select {
		case err := <-chanErr:
			return response, err
		case response := <-chanResp:
			return response, nil
		default:
			if waktu.Add(2 * time.Second).Before(time.Now()) {
				return response, errors.WithCode(code.GetAvatarTimeout, "Error timeout get avatar !")
			}
		}
	}

}

func (service userService) MyListGroups(_ context.Context) (response model.MyListGroupsResponse, err error) {
	whatsapp.MustLogin(service.waCli)

	groups, err := service.waCli.GetJoinedGroups()
	if err != nil {
		return
	}
	fmt.Printf("%+v\n", groups)
	if groups != nil {
		for _, group := range groups {
			response.Data = append(response.Data, *group)
		}
	}
	return response, nil
}

func (service userService) MyPrivacySetting(_ context.Context) (response model.MyPrivacySettingResponse, err error) {
	whatsapp.MustLogin(service.waCli)

	resp, err := service.waCli.TryFetchPrivacySettings(true)
	if err != nil {
		return
	}

	response.GroupAdd = string(resp.GroupAdd)
	response.Status = string(resp.Status)
	response.ReadReceipts = string(resp.ReadReceipts)
	response.Profile = string(resp.Profile)
	return response, nil
}

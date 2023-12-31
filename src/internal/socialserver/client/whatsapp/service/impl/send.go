package services

import (
	"context"
	"fmt"
	//fiberUtils "github.com/gofiber/fiber/v2/utils"
	//	"github.com/h2non/bimg"
	//	"github.com/marmotedu/errors"
	"github.com/valyala/fasthttp"
	//	"go-socialapp/internal/pkg/code"
	"go-socialapp/internal/pkg/third-party/whatsapp"
	"go-socialapp/internal/pkg/util"
	"go-socialapp/internal/socialserver/client/whatsapp/model"
	"go-socialapp/internal/socialserver/client/whatsapp/service/impl/validations"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"google.golang.org/protobuf/proto"
	"net/http"
	"os"
	//"os/exec"
)

type serviceSend struct {
	waClient   *whatsapp.WaClient
	appService *serviceApp
}

type metadata struct {
	Name    string
	Content string
}

//var sendSrv *serviceSend
//
//func GetSendService(waCli *whatsmeow.Client, appService interface{}) *serviceSend {
//	if sendSrv != nil {
//		return sendSrv
//	}
//	sendSrv = newSendService(waCli, appService)
//	return sendSrv
//}

func NewSendService(waClient *whatsapp.WaClient, appService interface{}) *serviceSend {
	return &serviceSend{
		waClient:   waClient,
		appService: appService.(*serviceApp),
	}
}

func (service serviceSend) SendText(ctx context.Context, request model.MessageRequest) (response model.MessageResponse, err error) {
	err = validations.ValidateSendMessage(ctx, request)
	if err != nil {
		return response, err
	}
	dataWaRecipient, err := service.waClient.ValidateJidWithLogin(request.Phone)
	if err != nil {
		return response, err
	}

	// Send message
	msg := &waProto.Message{Conversation: proto.String(request.Message)}

	// Reply message
	if request.ReplyMessageID != nil && *request.ReplyMessageID != "" {
		participantJID := dataWaRecipient.String()
		if len(*request.ReplyMessageID) < 28 {
			firstDevice, err := service.appService.FirstDevice(ctx)
			if err != nil {
				return response, err
			}
			participantJID = firstDevice.Device
		}

		msg = &waProto.Message{
			ExtendedTextMessage: &waProto.ExtendedTextMessage{
				Text: proto.String(request.Message),
				ContextInfo: &waProto.ContextInfo{
					StanzaId:    request.ReplyMessageID,
					Participant: proto.String(participantJID),
					QuotedMessage: &waProto.Message{
						Conversation: proto.String(request.Message),
					},
				},
			},
		}
	}

	ts, err := service.waClient.WaCli.SendMessage(ctx, dataWaRecipient, msg)
	if err != nil {
		return response, err
	}

	response.MessageID = ts.ID
	response.Status = fmt.Sprintf("Message sent to %s (server timestamp: %s)", request.Phone, ts)
	return response, nil
}

func (service serviceSend) SendImage(ctx context.Context, request model.ImageRequest) (response model.ImageResponse, err error) {
	//err = validations.ValidateSendImage(ctx, network)
	//if err != nil {
	//	return response, err
	//}
	//dataWaRecipient, err := whatsapp.ValidateJidWithLogin(service.waCli, network.Phone)
	//if err != nil {
	//	return response, err
	//}
	//
	//var (
	//	imagePath      string
	//	imageThumbnail string
	//	deletedItems   []string
	//)
	//
	//// Save image to server
	//oriImagePath := fmt.Sprintf("%s/%s", whatsapp.PathSendItems, network.Image.Filename)
	//err = fasthttp.SaveMultipartFile(network.Image, oriImagePath)
	//if err != nil {
	//	return response, err
	//}
	//deletedItems = append(deletedItems, oriImagePath)
	//
	//// Generate thumbnail with smalled image
	//openThumbnailBuffer, err := bimg.Read(oriImagePath)
	//imageThumbnail = fmt.Sprintf("%s/thumbnails-%s", whatsapp.PathSendItems, network.Image.Filename)
	//thumbnailImage, err := bimg.NewImage(openThumbnailBuffer).Process(bimg.Options{Quality: 90, Width: 100, Embed: true})
	//if err != nil {
	//	return response, err
	//}
	//err = bimg.Write(imageThumbnail, thumbnailImage)
	//if err != nil {
	//	return response, err
	//}
	//deletedItems = append(deletedItems, imageThumbnail)
	//
	//if network.Compress {
	//	// Resize image
	//	openImageBuffer, err := bimg.Read(oriImagePath)
	//	newImage, err := bimg.NewImage(openImageBuffer).Process(bimg.Options{Quality: 90, Width: 600, Embed: true})
	//	if err != nil {
	//		return response, err
	//	}
	//
	//	newImagePath := fmt.Sprintf("%s/new-%s", whatsapp.PathSendItems, network.Image.Filename)
	//	err = bimg.Write(newImagePath, newImage)
	//	if err != nil {
	//		return response, err
	//	}
	//	deletedItems = append(deletedItems, newImagePath)
	//	imagePath = newImagePath
	//} else {
	//	imagePath = oriImagePath
	//}
	//
	//// Send to WA server
	//dataWaCaption := network.Caption
	//dataWaImage, err := os.ReadFile(imagePath)
	//if err != nil {
	//	return response, err
	//}
	//uploadedImage, err := service.waCli.Upload(context.Background(), dataWaImage, whatsmeow.MediaImage)
	//if err != nil {
	//	fmt.Printf("Failed to upload file: %v", err)
	//	return response, err
	//}
	//dataWaThumbnail, err := os.ReadFile(imageThumbnail)
	//
	//msg := &waProto.Message{ImageMessage: &waProto.ImageMessage{
	//	JpegThumbnail: dataWaThumbnail,
	//	Caption:       proto.String(dataWaCaption),
	//	Url:           proto.String(uploadedImage.URL),
	//	DirectPath:    proto.String(uploadedImage.DirectPath),
	//	MediaKey:      uploadedImage.MediaKey,
	//	Mimetype:      proto.String(http.DetectContentType(dataWaImage)),
	//	FileEncSha256: uploadedImage.FileEncSHA256,
	//	FileSha256:    uploadedImage.FileSHA256,
	//	FileLength:    proto.Uint64(uint64(len(dataWaImage))),
	//	ViewOnce:      proto.Bool(network.ViewOnce),
	//}}
	//ts, err := service.waCli.SendMessage(ctx, dataWaRecipient, msg)
	//go func() {
	//	errDelete := utils.RemoveFile(0, deletedItems...)
	//	if errDelete != nil {
	//		fmt.Println("error when deleting picture: ", errDelete)
	//	}
	//}()
	//if err != nil {
	//	return response, err
	//}
	//
	//response.MessageID = ts.ID
	//response.Status = fmt.Sprintf("Message sent to %s (server timestamp: %s)", network.Phone, ts)
	return response, nil
}

func (service serviceSend) SendFile(ctx context.Context, request model.FileRequest) (response model.FileResponse, err error) {
	err = validations.ValidateSendFile(ctx, request)
	if err != nil {
		return response, err
	}
	dataWaRecipient, err := service.waClient.ValidateJidWithLogin(request.Phone)
	if err != nil {
		return response, err
	}

	oriFilePath := fmt.Sprintf("%s/%s", whatsapp.PathSendItems, request.File.Filename)
	err = fasthttp.SaveMultipartFile(request.File, oriFilePath)
	if err != nil {
		return response, err
	}

	// Send to WA server
	dataWaFile, err := os.ReadFile(oriFilePath)
	if err != nil {
		return response, err
	}
	uploadedFile, err := service.waClient.WaCli.Upload(context.Background(), dataWaFile, whatsmeow.MediaDocument)
	if err != nil {
		fmt.Printf("Failed to upload file: %v", err)
		return response, err
	}

	msg := &waProto.Message{DocumentMessage: &waProto.DocumentMessage{
		Url:           proto.String(uploadedFile.URL),
		Mimetype:      proto.String(http.DetectContentType(dataWaFile)),
		Title:         proto.String(request.File.Filename),
		FileSha256:    uploadedFile.FileSHA256,
		FileLength:    proto.Uint64(uploadedFile.FileLength),
		MediaKey:      uploadedFile.MediaKey,
		FileName:      proto.String(request.File.Filename),
		FileEncSha256: uploadedFile.FileEncSHA256,
		DirectPath:    proto.String(uploadedFile.DirectPath),
		Caption:       proto.String(request.Caption),
	}}
	ts, err := service.waClient.WaCli.SendMessage(ctx, dataWaRecipient, msg)
	go func() {
		errDelete := utils.RemoveFile(0, oriFilePath)
		if errDelete != nil {
			fmt.Println(errDelete)
		}
	}()
	if err != nil {
		return response, err
	}

	response.MessageID = ts.ID
	response.Status = fmt.Sprintf("Document sent to %s (server timestamp: %s)", request.Phone, ts)
	return response, nil
}

func (service serviceSend) SendVideo(ctx context.Context, request model.VideoRequest) (response model.VideoResponse, err error) {
	//err = validations.ValidateSendVideo(ctx, network)
	//if err != nil {
	//	return response, err
	//}
	//dataWaRecipient, err := whatsapp.ValidateJidWithLogin(service.waCli, network.Phone)
	//if err != nil {
	//	return response, err
	//}
	//
	//var (
	//	videoPath      string
	//	videoThumbnail string
	//	deletedItems   []string
	//)
	//
	//generateUUID := fiberUtils.UUIDv4()
	//// Save video to server
	//oriVideoPath := fmt.Sprintf("%s/%s", whatsapp.PathSendItems, generateUUID+network.Video.Filename)
	//err = fasthttp.SaveMultipartFile(network.Video, oriVideoPath)
	//if err != nil {
	//	return response, errors.WithCode(code.FailedStoreVideoInserver, fmt.Sprintf("failed to store video in server %v", err))
	//}
	//
	//// Get thumbnail video with ffmpeg
	//thumbnailVideoPath := fmt.Sprintf("%s/%s", whatsapp.PathSendItems, generateUUID+".png")
	//cmdThumbnail := exec.Command("ffmpeg", "-i", oriVideoPath, "-ss", "00:00:01.000", "-vframes", "1", thumbnailVideoPath)
	//err = cmdThumbnail.Run()
	//if err != nil {
	//	return response, errors.WithCode(code.FailedCreateThumbnail, fmt.Sprintf("failed to create thumbnail %v", err))
	//}
	//
	//// Resize Thumbnail
	//openImageBuffer, err := bimg.Read(thumbnailVideoPath)
	//resize, err := bimg.NewImage(openImageBuffer).Process(bimg.Options{Quality: 90, Width: 600, Embed: true})
	//if err != nil {
	//	return response, errors.WithCode(code.FailedResizeThumbail, fmt.Sprintf("failed to resize thumbnail %v", err))
	//}
	//thumbnailResizeVideoPath := fmt.Sprintf("%s/%s", whatsapp.PathSendItems, generateUUID+"_resize.png")
	//err = bimg.Write(thumbnailResizeVideoPath, resize)
	//if err != nil {
	//	return response, errors.WithCode(code.FailedCreateImageThumbnail, fmt.Sprintf("failed to create image thumbnail %v", err))
	//}
	//
	//deletedItems = append(deletedItems, thumbnailVideoPath)
	//deletedItems = append(deletedItems, thumbnailResizeVideoPath)
	//videoThumbnail = thumbnailResizeVideoPath
	//
	//if network.Compress {
	//	compresVideoPath := fmt.Sprintf("%s/%s", whatsapp.PathSendItems, generateUUID+".mp4")
	//	// Compress video with ffmpeg
	//	cmdCompress := exec.Command("ffmpeg", "-i", oriVideoPath, "-strict", "-2", compresVideoPath)
	//	err = cmdCompress.Run()
	//	if err != nil {
	//		return response, errors.WithCode(code.FailedCompressVideo, "failed to compress video")
	//	}
	//
	//	videoPath = compresVideoPath
	//	deletedItems = append(deletedItems, compresVideoPath)
	//} else {
	//	videoPath = oriVideoPath
	//	deletedItems = append(deletedItems, oriVideoPath)
	//}
	//
	////Send to WA server
	//dataWaVideo, err := os.ReadFile(videoPath)
	//if err != nil {
	//	return response, err
	//}
	//uploaded, err := service.waCli.Upload(context.Background(), dataWaVideo, whatsmeow.MediaVideo)
	//if err != nil {
	//	return response, errors.WithCode(code.FailedUploadFile, fmt.Sprintf("Failed to upload file: %v", err))
	//}
	//dataWaThumbnail, err := os.ReadFile(videoThumbnail)
	//if err != nil {
	//	return response, err
	//}
	//
	//msg := &waProto.Message{VideoMessage: &waProto.VideoMessage{
	//	Url:                 proto.String(uploaded.URL),
	//	Mimetype:            proto.String(http.DetectContentType(dataWaVideo)),
	//	Caption:             proto.String(network.Caption),
	//	FileLength:          proto.Uint64(uploaded.FileLength),
	//	FileSha256:          uploaded.FileSHA256,
	//	FileEncSha256:       uploaded.FileEncSHA256,
	//	MediaKey:            uploaded.MediaKey,
	//	DirectPath:          proto.String(uploaded.DirectPath),
	//	ViewOnce:            proto.Bool(network.ViewOnce),
	//	JpegThumbnail:       dataWaThumbnail,
	//	ThumbnailEncSha256:  dataWaThumbnail,
	//	ThumbnailSha256:     dataWaThumbnail,
	//	ThumbnailDirectPath: proto.String(uploaded.DirectPath),
	//}}
	//ts, err := service.waCli.SendMessage(ctx, dataWaRecipient, msg)
	//go func() {
	//	errDelete := utils.RemoveFile(1, deletedItems...)
	//	if errDelete != nil {
	//		fmt.Println(errDelete)
	//	}
	//}()
	//if err != nil {
	//	return response, err
	//}
	//
	//response.MessageID = ts.ID
	//response.Status = fmt.Sprintf("Video sent to %s (server timestamp: %s)", network.Phone, ts)
	return response, nil
}

func (service serviceSend) SendContact(ctx context.Context, request model.ContactRequest) (response model.ContactResponse, err error) {
	err = validations.ValidateSendContact(ctx, request)
	if err != nil {
		return response, err
	}
	dataWaRecipient, err := service.waClient.ValidateJidWithLogin(request.Phone)
	if err != nil {
		return response, err
	}

	msgVCard := fmt.Sprintf("BEGIN:VCARD\nVERSION:3.0\nN:;%v;;;\nFN:%v\nTEL;type=CELL;waid=%v:+%v\nEND:VCARD",
		request.ContactName, request.ContactName, request.ContactPhone, request.ContactPhone)
	msg := &waProto.Message{ContactMessage: &waProto.ContactMessage{
		DisplayName: proto.String(request.ContactName),
		Vcard:       proto.String(msgVCard),
	}}
	ts, err := service.waClient.WaCli.SendMessage(ctx, dataWaRecipient, msg)
	if err != nil {
		return response, err
	}

	response.MessageID = ts.ID
	response.Status = fmt.Sprintf("Contact sent to %s (server timestamp: %s)", request.Phone, ts)
	return response, nil
}

func (service serviceSend) SendLink(ctx context.Context, request model.LinkRequest) (response model.LinkResponse, err error) {
	err = validations.ValidateSendLink(ctx, request)
	if err != nil {
		return response, err
	}
	dataWaRecipient, err := service.waClient.ValidateJidWithLogin(request.Phone)
	if err != nil {
		return response, err
	}

	getMetaDataFromURL := utils.GetMetaDataFromURL(request.Link)

	msg := &waProto.Message{ExtendedTextMessage: &waProto.ExtendedTextMessage{
		Text:         proto.String(fmt.Sprintf("%s\n%s", request.Caption, request.Link)),
		Title:        proto.String(getMetaDataFromURL.Title),
		CanonicalUrl: proto.String(request.Link),
		MatchedText:  proto.String(request.Link),
		Description:  proto.String(getMetaDataFromURL.Description),
	}}
	ts, err := service.waClient.WaCli.SendMessage(ctx, dataWaRecipient, msg)
	if err != nil {
		return response, err
	}

	response.MessageID = ts.ID
	response.Status = fmt.Sprintf("Link sent to %s (server timestamp: %s)", request.Phone, ts)
	return response, nil
}

func (service serviceSend) SendLocation(ctx context.Context, request model.LocationRequest) (response model.LocationResponse, err error) {
	err = validations.ValidateSendLocation(ctx, request)
	if err != nil {
		return response, err
	}
	dataWaRecipient, err := service.waClient.ValidateJidWithLogin(request.Phone)
	if err != nil {
		return response, err
	}

	// Compose WhatsApp Proto
	msg := &waProto.Message{
		LocationMessage: &waProto.LocationMessage{
			DegreesLatitude:  proto.Float64(utils.StrToFloat64(request.Latitude)),
			DegreesLongitude: proto.Float64(utils.StrToFloat64(request.Longitude)),
		},
	}

	// Send WhatsApp Message Proto
	ts, err := service.waClient.WaCli.SendMessage(ctx, dataWaRecipient, msg)
	if err != nil {
		return response, err
	}

	response.MessageID = ts.ID
	response.Status = fmt.Sprintf("Send location success %s (server timestamp: %s)", request.Phone, ts)
	return response, nil
}

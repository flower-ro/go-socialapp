package validations

import (
	"context"
	"fmt"
	"github.com/dustin/go-humanize"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/marmotedu/errors"
	"go-socialapp/internal/pkg/code"
	"go-socialapp/internal/pkg/third-party/whatsapp"
	"go-socialapp/internal/socialserver/client/whatsapp/model"
)

func ValidateSendMessage(ctx context.Context, request model.MessageRequest) error {
	err := validation.ValidateStructWithContext(ctx, &request,
		validation.Field(&request.Phone, validation.Required),
		validation.Field(&request.Message, validation.Required),
	)

	if err != nil {
		return errors.WithCode(code.ValidationError, err.Error())
	}
	return nil
}

func ValidateSendImage(ctx context.Context, request model.ImageRequest) error {
	err := validation.ValidateStructWithContext(ctx, &request,
		validation.Field(&request.Phone, validation.Required),
		validation.Field(&request.Image, validation.Required),
	)

	if err != nil {
		return errors.WithCode(code.ValidationError, err.Error())
	}

	availableMimes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
	}

	if !availableMimes[request.Image.Header.Get("Content-Type")] {
		return errors.WithCode(code.ValidationError, "your image is not allowed. please use jpg/jpeg/png")
	}

	return nil
}

func ValidateSendFile(ctx context.Context, request model.FileRequest) error {
	err := validation.ValidateStructWithContext(ctx, &request,
		validation.Field(&request.Phone, validation.Required),
		validation.Field(&request.File, validation.Required),
	)

	if err != nil {
		return errors.WithCode(code.ValidationError, err.Error())
	}

	if request.File.Size > whatsapp.WhatsappSettingMaxFileSize { // 10MB
		maxSizeString := humanize.Bytes(uint64(whatsapp.WhatsappSettingMaxFileSize))
		return errors.WithCode(code.ValidationError, fmt.Sprintf("max file upload is %s, please upload in cloud and send via text if your file is higher than %s", maxSizeString, maxSizeString))
	}

	return nil
}

func ValidateSendVideo(ctx context.Context, request model.VideoRequest) error {
	err := validation.ValidateStructWithContext(ctx, &request,
		validation.Field(&request.Phone, validation.Required),
		validation.Field(&request.Video, validation.Required),
	)

	if err != nil {
		return errors.WithCode(code.ValidationError, err.Error())
	}

	availableMimes := map[string]bool{
		"video/mp4":        true,
		"video/x-matroska": true,
		"video/avi":        true,
	}

	if !availableMimes[request.Video.Header.Get("Content-Type")] {
		return errors.WithCode(code.ValidationError, "your video type is not allowed. please use mp4/mkv/avi")
	}

	if request.Video.Size > whatsapp.WhatsappSettingMaxVideoSize { // 30MB
		maxSizeString := humanize.Bytes(uint64(whatsapp.WhatsappSettingMaxVideoSize))
		return errors.WithCode(code.ValidationError, fmt.Sprintf("max video upload is %s, please upload in cloud and send via text if your file is higher than %s", maxSizeString, maxSizeString))
	}

	return nil
}

func ValidateSendContact(ctx context.Context, request model.ContactRequest) error {
	err := validation.ValidateStructWithContext(ctx, &request,
		validation.Field(&request.Phone, validation.Required),
		validation.Field(&request.ContactPhone, validation.Required),
		validation.Field(&request.ContactName, validation.Required),
	)

	if err != nil {
		return errors.WithCode(code.ValidationError, err.Error())
	}

	return nil
}

func ValidateSendLink(ctx context.Context, request model.LinkRequest) error {
	err := validation.ValidateStructWithContext(ctx, &request,
		validation.Field(&request.Phone, validation.Required),
		validation.Field(&request.Link, validation.Required, is.URL),
		validation.Field(&request.Caption, validation.Required),
	)

	if err != nil {
		return errors.WithCode(code.ValidationError, err.Error())
	}

	return nil
}

func ValidateSendLocation(ctx context.Context, request model.LocationRequest) error {
	err := validation.ValidateStructWithContext(ctx, &request,
		validation.Field(&request.Phone, validation.Required),
		validation.Field(&request.Latitude, validation.Required, is.Latitude),
		validation.Field(&request.Longitude, validation.Required, is.Longitude),
	)

	if err != nil {
		return errors.WithCode(code.ValidationError, err.Error())
	}

	return nil
}

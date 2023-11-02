package validations

import (
	"context"
	"github.com/marmotedu/errors"
	"go-socialapp/internal/pkg/code"
	"go-socialapp/internal/socialserver/client/whatsapp/model"

	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"testing"
)

func TestValidateSendMessage(t *testing.T) {
	type args struct {
		request model.MessageRequest
	}
	tests := []struct {
		name string
		args args
		err  any
	}{
		{
			name: "should success with phone and message",
			args: args{request: model.MessageRequest{
				Phone:   "1728937129312@s.whatsapp.net",
				Message: "Hello this is testing",
			}},
			err: nil,
		},
		{
			name: "should error with empty phone",
			args: args{request: model.MessageRequest{
				Phone:   "",
				Message: "Hello this is testing",
			}},
			err: errors.WithCode(code.ValidationError, "phone: cannot be blank."),
		},
		{
			name: "should error with empty message",
			args: args{request: model.MessageRequest{
				Phone:   "1728937129312@s.whatsapp.net",
				Message: "",
			}},
			err: errors.WithCode(code.ValidationError, "message: cannot be blank."),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSendMessage(context.Background(), tt.args.request)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestValidateSendImage(t *testing.T) {
	image := &multipart.FileHeader{
		Filename: "sample-image.png",
		Size:     100,
		Header:   map[string][]string{"Content-Type": {"image/png"}},
	}

	type args struct {
		request model.ImageRequest
	}
	tests := []struct {
		name string
		args args
		err  any
	}{
		{
			name: "should success with normal condition",
			args: args{request: model.ImageRequest{
				Phone:   "1728937129312@s.whatsapp.net",
				Caption: "Hello this is testing",
				Image:   image,
			}},
			err: nil,
		},
		{
			name: "should error with empty phone",
			args: args{request: model.ImageRequest{
				Phone: "",
				Image: image,
			}},
			err: errors.WithCode(code.ValidationError, "phone: cannot be blank."),
		},
		{
			name: "should error with empty image",
			args: args{request: model.ImageRequest{
				Phone: "1728937129312@s.whatsapp.net",
				Image: nil,
			}},
			err: errors.WithCode(code.ValidationError, "image: cannot be blank."),
		},
		{
			name: "should error with invalid image type",
			args: args{request: model.ImageRequest{
				Phone: "1728937129312@s.whatsapp.net",
				Image: &multipart.FileHeader{
					Filename: "sample-image.pdf",
					Size:     100,
					Header:   map[string][]string{"Content-Type": {"application/pdf"}},
				},
			}},
			err: errors.WithCode(code.ValidationError, "your image is not allowed. please use jpg/jpeg/png"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSendImage(context.Background(), tt.args.request)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestValidateSendFile(t *testing.T) {
	file := &multipart.FileHeader{
		Filename: "sample-image.png",
		Size:     100,
		Header:   map[string][]string{"Content-Type": {"image/png"}},
	}

	type args struct {
		request model.FileRequest
	}
	tests := []struct {
		name string
		args args
		err  any
	}{
		{
			name: "should success with normal condition",
			args: args{request: model.FileRequest{
				Phone: "1728937129312@s.whatsapp.net",
				File:  file,
			}},
			err: nil,
		},
		{
			name: "should error with empty phone",
			args: args{request: model.FileRequest{
				Phone: "",
				File:  file,
			}},
			err: errors.WithCode(code.ValidationError, "phone: cannot be blank."),
		},
		{
			name: "should error with empty file",
			args: args{request: model.FileRequest{
				Phone: "1728937129312@s.whatsapp.net",
				File:  nil,
			}},
			err: errors.WithCode(code.ValidationError, "file: cannot be blank."),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSendFile(context.Background(), tt.args.request)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestValidateSendVideo(t *testing.T) {
	file := &multipart.FileHeader{
		Filename: "sample-video.mp4",
		Size:     100,
		Header:   map[string][]string{"Content-Type": {"video/mp4"}},
	}

	type args struct {
		request model.VideoRequest
	}
	tests := []struct {
		name string
		args args
		err  any
	}{
		{
			name: "should success with normal condition",
			args: args{request: model.VideoRequest{
				Phone:    "1728937129312@s.whatsapp.net",
				Caption:  "simple caption",
				Video:    file,
				ViewOnce: false,
				Compress: false,
			}},
			err: nil,
		},
		{
			name: "should error with empty phone",
			args: args{request: model.VideoRequest{
				Phone:    "",
				Caption:  "simple caption",
				Video:    file,
				ViewOnce: false,
				Compress: false,
			}},
			err: errors.WithCode(code.ValidationError, "phone: cannot be blank."),
		},
		{
			name: "should error with empty video",
			args: args{request: model.VideoRequest{
				Phone:    "1728937129312@s.whatsapp.net",
				Caption:  "simple caption",
				Video:    nil,
				ViewOnce: false,
				Compress: false,
			}},
			err: errors.WithCode(code.ValidationError, "video: cannot be blank."),
		},
		{
			name: "should error with invalid format video",
			args: args{request: model.VideoRequest{
				Phone:   "1728937129312@s.whatsapp.net",
				Caption: "simple caption",
				Video: func() *multipart.FileHeader {
					return &multipart.FileHeader{
						Filename: "sample-video.jpg",
						Size:     100,
						Header:   map[string][]string{"Content-Type": {"image/png"}},
					}
				}(),
				ViewOnce: false,
				Compress: false,
			}},
			err: errors.WithCode(code.ValidationError, "your video type is not allowed. please use mp4/mkv/avi"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSendVideo(context.Background(), tt.args.request)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestValidateSendLink(t *testing.T) {
	type args struct {
		request model.LinkRequest
	}
	tests := []struct {
		name string
		args args
		err  any
	}{
		{
			name: "should success normal condition",
			args: args{request: model.LinkRequest{
				Phone:   "1728937129312@s.whatsapp.net",
				Caption: "description",
				Link:    "https://google.com",
			}},
			err: nil,
		},
		{
			name: "should error with empty phone",
			args: args{request: model.LinkRequest{
				Phone:   "",
				Caption: "description",
				Link:    "https://google.com",
			}},
			err: errors.WithCode(code.ValidationError, "phone: cannot be blank."),
		},
		{
			name: "should error with empty caption",
			args: args{request: model.LinkRequest{
				Phone:   "1728937129312@s.whatsapp.net",
				Caption: "",
				Link:    "https://google.com",
			}},
			err: errors.WithCode(code.ValidationError, "caption: cannot be blank."),
		},
		{
			name: "should error with empty link",
			args: args{request: model.LinkRequest{
				Phone:   "1728937129312@s.whatsapp.net",
				Caption: "description",
				Link:    "",
			}},
			err: errors.WithCode(code.ValidationError, "link: cannot be blank."),
		},
		{
			name: "should error with invalid link",
			args: args{request: model.LinkRequest{
				Phone:   "1728937129312@s.whatsapp.net",
				Caption: "description",
				Link:    "googlecom",
			}},
			err: errors.WithCode(code.ValidationError, "link: must be a valid URL."),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSendLink(context.Background(), tt.args.request)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestValidateRevokeMessage(t *testing.T) {
	type args struct {
		request model.RevokeRequest
	}
	tests := []struct {
		name string
		args args
		err  any
	}{
		{
			name: "should success normal condition",
			args: args{request: model.RevokeRequest{
				Phone:     "1728937129312@s.whatsapp.net",
				MessageID: "1382901271239781",
			}},
			err: nil,
		},
		{
			name: "should error with empty phone",
			args: args{request: model.RevokeRequest{
				Phone:     "",
				MessageID: "1382901271239781",
			}},
			err: errors.WithCode(code.ValidationError, "phone: cannot be blank."),
		},
		{
			name: "should error with empty message id",
			args: args{request: model.RevokeRequest{
				Phone:     "1728937129312@s.whatsapp.net",
				MessageID: "",
			}},
			err: errors.WithCode(code.ValidationError, "message_id: cannot be blank."),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRevokeMessage(context.Background(), tt.args.request)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestValidateUpdateMessage(t *testing.T) {
	type args struct {
		request model.UpdateMessageRequest
	}
	tests := []struct {
		name string
		args args
		err  any
	}{
		{
			name: "should success normal condition",
			args: args{request: model.UpdateMessageRequest{
				MessageID: "1382901271239781",
				Message:   "some update message",
				Phone:     "1728937129312@s.whatsapp.net",
			}},
			err: nil,
		},
		{
			name: "should error with empty phone",
			args: args{request: model.UpdateMessageRequest{
				MessageID: "1382901271239781",
				Message:   "some update message",
				Phone:     "",
			}},
			err: errors.WithCode(code.ValidationError, "phone: cannot be blank."),
		},
		{
			name: "should error with empty message id",
			args: args{request: model.UpdateMessageRequest{
				MessageID: "",
				Message:   "some update message",
				Phone:     "1728937129312@s.whatsapp.net",
			}},
			err: errors.WithCode(code.ValidationError, "message_id: cannot be blank."),
		},
		{
			name: "should error with empty message update",
			args: args{request: model.UpdateMessageRequest{
				MessageID: "1382901271239781",
				Message:   "",
				Phone:     "1728937129312@s.whatsapp.net",
			}},
			err: errors.WithCode(code.ValidationError, "message: cannot be blank."),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUpdateMessage(context.Background(), tt.args.request)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestValidateSendContact(t *testing.T) {
	type args struct {
		request model.ContactRequest
	}
	tests := []struct {
		name string
		args args
		err  any
	}{
		{
			name: "should success normal condition",
			args: args{request: model.ContactRequest{
				Phone:        "1728937129312@s.whatsapp.net",
				ContactName:  "Aldino",
				ContactPhone: "62788712738123",
			}},
			err: nil,
		},
		{
			name: "should error with empty phone",
			args: args{request: model.ContactRequest{
				Phone:        "",
				ContactName:  "Aldino",
				ContactPhone: "62788712738123",
			}},
			err: errors.WithCode(code.ValidationError, "phone: cannot be blank."),
		},
		{
			name: "should error with empty contact name",
			args: args{request: model.ContactRequest{
				Phone:        "1728937129312@s.whatsapp.net",
				ContactName:  "",
				ContactPhone: "62788712738123",
			}},
			err: errors.WithCode(code.ValidationError, "contact_name: cannot be blank."),
		},
		{
			name: "should error with empty contact phone",
			args: args{request: model.ContactRequest{
				Phone:        "1728937129312@s.whatsapp.net",
				ContactName:  "Aldino",
				ContactPhone: "",
			}},
			err: errors.WithCode(code.ValidationError, "contact_phone: cannot be blank."),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSendContact(context.Background(), tt.args.request)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestValidateSendLocation(t *testing.T) {
	type args struct {
		request model.LocationRequest
	}
	tests := []struct {
		name string
		args args
		err  any
	}{
		{
			name: "should success normal condition",
			args: args{request: model.LocationRequest{
				Phone:     "1728937129312@s.whatsapp.net",
				Latitude:  "-7.797068",
				Longitude: "110.370529",
			}},
			err: nil,
		},
		{
			name: "should error with empty phone",
			args: args{request: model.LocationRequest{
				Phone:     "",
				Latitude:  "-7.797068",
				Longitude: "110.370529",
			}},
			err: errors.WithCode(code.ValidationError, "phone: cannot be blank."),
		},
		{
			name: "should error with empty latitude",
			args: args{request: model.LocationRequest{
				Phone:     "1728937129312@s.whatsapp.net",
				Latitude:  "",
				Longitude: "110.370529",
			}},
			err: errors.WithCode(code.ValidationError, "latitude: cannot be blank."),
		},
		{
			name: "should error with empty longitude",
			args: args{request: model.LocationRequest{
				Phone:     "1728937129312@s.whatsapp.net",
				Latitude:  "-7.797068",
				Longitude: "",
			}},
			err: errors.WithCode(code.ValidationError, "longitude: cannot be blank."),
		},
		{
			name: "should error with invalid latitude",
			args: args{request: model.LocationRequest{
				Phone:     "1728937129312@s.whatsapp.net",
				Latitude:  "ABCDEF",
				Longitude: "110.370529",
			}},
			err: errors.WithCode(code.ValidationError, "latitude: must be a valid latitude."),
		},
		{
			name: "should error with invalid latitude",
			args: args{request: model.LocationRequest{
				Phone:     "1728937129312@s.whatsapp.net",
				Latitude:  "-7.797068",
				Longitude: "ABCDEF",
			}},
			err: errors.WithCode(code.ValidationError, "longitude: must be a valid longitude."),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSendLocation(context.Background(), tt.args.request)
			assert.Equal(t, tt.err, err)
		})
	}
}

package model

import "mime/multipart"

// contact
type ContactRequest struct {
	Phone        string `json:"phone" form:"phone"`
	ContactName  string `json:"contact_name" form:"contact_name"`
	ContactPhone string `json:"contact_phone" form:"contact_phone"`
}
type ContactResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
}

// file
type FileRequest struct {
	Phone   string                `json:"phone" form:"phone"`
	File    *multipart.FileHeader `json:"file" form:"file"`
	Caption string                `json:"caption" form:"caption"`
}
type FileResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
}

// image
type ImageRequest struct {
	Phone    string                `json:"phone" form:"phone"`
	Caption  string                `json:"caption" form:"caption"`
	Image    *multipart.FileHeader `json:"image" form:"image"`
	ViewOnce bool                  `json:"view_once" form:"view_once"`
	Compress bool                  `json:"compress"`
}
type ImageResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
}

// link
type LinkRequest struct {
	Phone   string `json:"phone" form:"phone"`
	Caption string `json:"caption"`
	Link    string `json:"link"`
}
type LinkResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
}

// location
type LocationRequest struct {
	Phone     string `json:"phone" form:"phone"`
	Latitude  string `json:"latitude" form:"latitude"`
	Longitude string `json:"longitude" form:"longitude"`
}
type LocationResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
}

// text
type MessageRequest struct {
	Phone          string  `json:"phone" form:"phone"`
	Message        string  `json:"message" form:"message"`
	ReplyMessageID *string `json:"reply_message_id" form:"reply_message_id"`
}
type MessageResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
}

// vedio
type VideoRequest struct {
	Phone    string                `json:"phone" form:"phone"`
	Caption  string                `json:"caption" form:"caption"`
	Video    *multipart.FileHeader `json:"video" form:"video"`
	ViewOnce bool                  `json:"view_once" form:"view_once"`
	Compress bool                  `json:"compress"`
}
type VideoResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
}

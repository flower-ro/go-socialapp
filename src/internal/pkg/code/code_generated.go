// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Code generated by "codegen -type=int /home/colin/workspace/golang/src/github.com/marmotedu/iam/internal/pkg/code"; DO NOT EDIT.

package code

// init register error codes defines in this source code to `github.com/marmotedu/errors`
func init() {
	register(ErrSuccess, 200, "OK")
	register(ErrUnknown, 500, "Internal server error")
	register(ErrBind, 400, "Error occurred while binding the request body to the struct")
	register(ErrValidation, 400, "Validation failed")
	register(ErrTokenInvalid, 401, "Token invalid")
	register(ErrPageNotFound, 404, "Page not found")
	register(ErrDatabase, 500, "Database error")
	register(ErrEncrypt, 401, "Error occurred while encrypting the account password")
	register(ErrSignatureInvalid, 401, "Signature is invalid")
	register(ErrExpired, 401, "Token expired")
	register(ErrInvalidAuthHeader, 401, "Invalid authorization header")
	register(ErrMissingHeader, 401, "The `Authorization` header was empty")
	register(ErrPasswordIncorrect, 401, "Password was incorrect")
	register(ErrPermissionDenied, 403, "Permission denied")
	register(ErrEncodingFailed, 500, "Encoding failed due to an error with the data")
	register(ErrDecodingFailed, 500, "Decoding failed due to an error with the data")
	register(ErrInvalidJSON, 500, "Data is not valid JSON")
	register(ErrEncodingJSON, 500, "JSON data could not be encoded")
	register(ErrDecodingJSON, 500, "JSON data could not be decoded")
	register(ErrInvalidYaml, 500, "Data is not valid Yaml")
	register(ErrEncodingYaml, 500, "Yaml data could not be encoded")
	register(ErrDecodingYaml, 500, "Yaml data could not be decoded")

	register(ErrConvertDto2Do, 500, "Dto to do Converting failed")
	register(ErrConvertDo2Dto, 500, "do to dto Converting failed")
	register(ErrJsonMarshal, 500, "Json marshal error")
	register(ErrJsonUnMarshal, 500, "Json unmarshal error")
	register(ErrCallService, 500, "service call err")
	register(ErrHandleServiceReturnResult, 500, "service call result handle err")
	register(ErrCacheSet, 500, "ErrCacheSet")
	register(ErrCacheNotFound, 500, "ErrCacheNotFound")
	register(ErrCacheGet, 500, "ErrCacheGet")

	register(ErrData, 500, "Data err")
	register(ErrRecordNotExisted, 500, "Record not existed.")

	register(ErrWebhook, 500, "WEBHOOK_ERROR.")
	register(ErrWaCLI, 500, "your WhatsApp CLI is invalid or empty.")
	register(ErrAlreadyLoggedIn, 500, "You already logged in.")
	register(ErrReconnect, 500, "Reconnect error.")
	register(ErrQrChannel, 500, "QR channel error.")
	register(ErrSessionSaved, 500, "Your session have been saved, please wait to connect 2 second and refresh again.")
	register(ClientNotInitialized, 500, "Whatsapp client is not initialized.")
	register(NotConnectServer, 500, "you are not connect to services server, please reconnect.")
	register(NotLoginServer, 500, "you are not login to services server, please login.")
	register(ErrInvalidJID, 500, "your JID is invalid.")
	register(GetAvatarTimeout, 500, "Error timeout get avatar.")
	register(ValidationError, 500, "validation err.")
	register(FailedStoreVideoInserver, 500, "failed to store video in server.")
	register(FailedCreateThumbnail, 500, "failed to create thumbnail.")
	register(FailedResizeThumbail, 500, "failed to resize thumbnail.")
	register(FailedCreateImageThumbnail, 500, "failed to create image thumbnail.")
	register(FailedCompressVideo, 500, "failed to compress video.")
	register(FailedUploadFile, 500, "Failed to upload file.")
}

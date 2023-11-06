package code

const (

	// ErrWebhook - 500: WEBHOOK_ERROR.
	ErrWebhook int = iota + 101001

	// ErrWaCLI - 500: your WhatsApp CLI is invalid or empty.
	ErrWaCLI
	// ErrAlreadyLoggedIn - 500: You already logged in.
	ErrAlreadyLoggedIn
	// ErrReconnect - 500: Reconnect error.
	ErrReconnect
	// ErrQrChannel - 500: QR channel error.
	ErrQrChannel
	// ErrSessionSaved - 500: Your session have been saved, please wait to connect 2 second and refresh again.
	ErrSessionSaved

	// ClientNotInitialized - 500: Whatsapp client is not initialized.
	ClientNotInitialized
	// NotConnectServer - 500: you are not connect to services server, please reconnect.
	NotConnectServer
	// NotLoginServer - 500: you are not login to services server, please login.
	NotLoginServer
	// ErrInvalidJID - 500: your JID is invalid.
	ErrInvalidJID
	// GetAvatarTimeout - 500: Error timeout get avatar.
	GetAvatarTimeout
	// ValidationError - 500: validation err.
	ValidationError

	// FailedStoreVideoInserver - 500: failed to store video in server.
	FailedStoreVideoInserver
	// FailedCreateThumbnail - 500: failed to create thumbnail.
	FailedCreateThumbnail
	// FailedResizeThumbail - 500: failed to resize thumbnail.
	FailedResizeThumbail
	// FailedCreateImageThumbnail - 500: failed to create image thumbnail .
	FailedCreateImageThumbnail
	// FailedCompressVideo - 500: failed to compress video.
	FailedCompressVideo
	// FailedUploadFile - 500: Failed to upload file.
	FailedUploadFile
	// FailedGetDevice - 500: Failed to get device.
	FailedGetDevice
	// FailedConnectSqlite3 - 500: Failed to connect to sqlite3 database.
	FailedConnectSqlite3
	// WaClientExistedInCache - 500: existed whatsapp waClient in cache.
	WaClientExistedInCache

	// DirNotExistedErr - 500: file dir is not existed.
	DirNotExistedErr

	// NotSureIsExisted - 500: file Or dir is not sure existed.
	NotSureIsExisted

	// FileIsExisted - 500: file is existed.
	FileIsExisted

	// FileIsNotExisted - 500: file is not existed.
	FileIsNotExisted

	// FileCreatedFail - 500: failed to create file.
	FileCreatedFail

	// NotLoginErr - 500: not login.
	NotLoginErr
)

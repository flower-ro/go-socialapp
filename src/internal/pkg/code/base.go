package code

//go:generate codegen -type=int
//go:generate codegen -type=int -doc -output ../../../docs/guide/zh-CN/api/error_code_generated.md

// Common: basic errors.
// Code must start with 1xxxxx.
const (
	// ErrSuccess - 200: OK.
	ErrSuccess int = iota + 100001

	// ErrUnknown - 500: Internal server error.
	ErrUnknown

	// ErrBind - 400: Error occurred while binding the network body to the struct.
	ErrBind

	// ErrValidation - 400: Validation failed.
	ErrValidation

	// ErrTokenInvalid - 401: Token invalid.
	ErrTokenInvalid

	// ErrPageNotFound - 404: Page not found.
	ErrPageNotFound
)

// common: database errors.
const (
	// ErrDatabase - 500: Database error.
	ErrDatabase int = iota + 100101
)

// common: authorization and authentication errors.
const (
	// ErrEncrypt - 401: Error occurred while encrypting the account password.
	ErrEncrypt int = iota + 100201

	// ErrSignatureInvalid - 401: Signature is invalid.
	ErrSignatureInvalid

	// ErrExpired - 401: Token expired.
	ErrExpired

	// ErrInvalidAuthHeader - 401: Invalid authorization header.
	ErrInvalidAuthHeader

	// ErrMissingHeader - 401: The `Authorization` header was empty.
	ErrMissingHeader

	// ErrPasswordIncorrect - 401: Password was incorrect.
	ErrPasswordIncorrect

	// PermissionDenied - 403: Permission denied.
	ErrPermissionDenied

	// ErrTLSGenerate - 403: Permission denied.
	ErrTLSGenerate
)

// common: encode/decode errors.
const (
	// ErrEncodingFailed - 500: Encoding failed due to an error with the data.
	ErrEncodingFailed int = iota + 100301

	// ErrDecodingFailed - 500: Decoding failed due to an error with the data.
	ErrDecodingFailed

	// ErrInvalidJSON - 500: Data is not valid JSON.
	ErrInvalidJSON

	// ErrEncodingJSON - 500: JSON data could not be encoded.
	ErrEncodingJSON

	// ErrDecodingJSON - 500: JSON data could not be decoded.
	ErrDecodingJSON

	// ErrInvalidYaml - 500: Data is not valid Yaml.
	ErrInvalidYaml

	// ErrEncodingYaml - 500: Yaml data could not be encoded.
	ErrEncodingYaml

	// ErrDecodingYaml - 500: Yaml data could not be decoded.
	ErrDecodingYaml
)

const (
	// ErrConvertDto2Do - 500: Dto to do Converting failed.
	ErrConvertDto2Do int = iota + 100401

	// ErrConvertDo2Dto - 500: do to dto Converting failed.
	ErrConvertDo2Dto

	// ErrJsonMarshal - 500: Json marshal error.
	ErrJsonMarshal

	// ErrJsonUnMarshal - 500: Json unmarshal error.
	ErrJsonUnMarshal

	// ErrCallService - 400: service call err.
	ErrCallService

	// ErrHandleServiceReturnResult - 400: service call result handle err.
	ErrHandleServiceReturnResult

	// ErrCacheSet - 500:  set cache error.
	ErrCacheSet

	// ErrCacheNotFound - 404: not found in cache.
	ErrCacheNotFound

	// ErrCacheGet - 500: get cache error.
	ErrCacheGet

	// ErrData - 500: Data err.
	ErrData
	// ErrRecordNotExisted - 500: Record not existed.
	ErrRecordNotExisted
)

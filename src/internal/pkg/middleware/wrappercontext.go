package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gostool/jsonconv"
	"github.com/marmotedu/errors"
	"github.com/marmotedu/iam/pkg/log"
	"go-socialapp/internal/pkg/code"
	"net/http"
)

const (
	JsonFormat       = "jsonFormat"
	DefaultErrorCode = 1000
	Authority        = "authority"
)

type ErrResponse struct {
	// Code defines the business error code.
	Code int `json:"code"`

	// Message contains the detail of this message.
	// This message is suitable to be exposed to external
	Message string `json:"message"`

	// Reference returns the reference document which maybe useful to solve this error.
	Reference string `json:"reference,omitempty"`
}

type WrapperContext struct {
	context    *gin.Context
	err        error
	original   bool
	httpStatus int
	resultData interface{}
	//errorMessage string
	//errorCode    int
	header map[string]string
}

func NewWrapperContext(c *gin.Context) *WrapperContext {
	var header = make(map[string]string)
	return &WrapperContext{
		context:    c,
		err:        nil,
		original:   false,
		httpStatus: http.StatusOK,
		resultData: nil,
		header:     header,
	}
}

type ApiHandlerFunc func(wc *WrapperContext)

func DealHanlder(h ApiHandlerFunc) func(c *gin.Context) {
	return func(c *gin.Context) {
		wc := NewWrapperContext(c)
		//
		h(wc)
		defer wc.Write()
	}
}

func (wc *WrapperContext) Context() *gin.Context {
	return wc.context
}

// 解析校验请求
func (wc *WrapperContext) ShouldBindJSON(i interface{}) bool {
	err := wc.context.ShouldBindJSON(i)
	if err != nil {
		wc.ErrorsWithCode(code.ErrBind, err.Error(), nil)
		//	wc.ParseParamError("body")
		return true
	}
	return false
}

func (wc *WrapperContext) ShouldBindQuery(i interface{}) bool {
	err := wc.context.ShouldBindQuery(i)
	if err != nil {
		wc.ErrorsWithCode(code.ErrBind, err.Error(), nil)
		//	wc.ParseParamError("body")
		return true
	}
	return false
}

func (wc *WrapperContext) QueryArray(key string) []string {
	return wc.context.QueryArray(key)
}

func (wc *WrapperContext) Param(param string) string {
	return wc.context.Param(param)
}

func (wc *WrapperContext) Query(param string) string {
	return wc.context.Query(param)
}

//func (wc *WrapperContext) ParseParamError(params ...string) {
//	var message = " parameter is invalid"
//	if len(params) > 0 {
//		message = strings.Join(params, ",") + message
//	}
//	wc.Errors(errors.New(message))
//}

//func (wc *WrapperContext) RequireToken() (string, bool) {
//	token := wc.context.GetHeader(models.Authorization)
//	if token == "" {
//		wc.err = errors.New("token is null")
//		return token, false
//	}
//	return token, true
//}

// 从请求里面取出参数
//
//	func (wc *WrapperContext) GetToken() string {
//		return wc.context.GetHeader(models.Authorization)
//	}
//
// func (wc *WrapperContext) GetRequestId() string { //RequestId前面塞进去
//
//		return wc.context.GetString(models.XRequestID)
//	}
func (wc *WrapperContext) GetAuthorities() string { //前面解析塞进去-这个好像没有用
	return wc.context.GetString(Authority)
}
func (wc *WrapperContext) GetUserName() string { //前面解析塞进去
	return wc.context.GetString(UsernameKey)
}

//func (wc *WrapperContext) GetRoles() string { //前面解析塞进去
//	return wc.context.GetString(models.Roles)
//}

// 返回
func (wc *WrapperContext) OriginalData(data []byte) {
	wc.original = true
	wc.resultData = data
}

func (wc *WrapperContext) Write() {
	if wc.header != nil {
		writeHeader(wc)
	}
	if wc.err == nil {
		success(wc)
	} else {
		fail(wc)
	}
}

func fail(wc *WrapperContext) {
	c := wc.context
	log.Errorf("%#+v", wc.err)
	coder := errors.ParseCoder(wc.err)
	c.JSON(coder.HTTPStatus(), ErrResponse{
		Code:      coder.Code(),
		Message:   coder.String(),
		Reference: coder.Reference(),
	})
}
func success(wc *WrapperContext) {
	c := wc.context
	format := wc.context.GetString(JsonFormat)
	switch format {
	case "camel":
		c.JSON(http.StatusOK, jsonconv.JsonCamelCase{
			Value: wc.resultData,
		})
	case "snake":
		c.JSON(http.StatusOK, jsonconv.JsonSnakeCase{
			Value: wc.resultData,
		})
	default:
		c.JSON(http.StatusOK, wc.resultData)
	}
}

func writeHeader(wc *WrapperContext) {
	c := wc.context
	for k, v := range wc.header {
		c.Header(k, v)
	}
}

func (wc *WrapperContext) ResultCamel() {
	wc.context.Set(JsonFormat, "snake")
}

//func (wc *WrapperContext) SetHttpStatus(status int) {
//	wc.httpStatus = status
//}

func (wc *WrapperContext) Success(data interface{}) {
	wc.resultData = data
}

// 分页返回
//func (wc *WrapperContext) SuccessPageInfo(startIndex, pageSize, total int, data interface{}) {
//	result := make(map[string]interface{})
//	result["pageInfo"] = &dto.PageInfo{
//		startIndex,
//		pageSize,
//		total,
//	}
//	result["content"] = data
//	wc.resultData = result
//}

// 数组返回-一般不需要
//func (wc *WrapperContext) SuccessList(data interface{}) {
//	result := make(map[string]interface{})
//	result["contents"] = data
//	wc.resultData = result
//}

// 工具方法
func (wc *WrapperContext) AddHeader(k, v string) {
	wc.header[k] = v
}

// error 相关方法
func (wc *WrapperContext) ErrorsWithCode(code int, format string, args ...interface{}) {
	wc.err = errors.WithCode(code, format, args)
}

func (wc *WrapperContext) Errors(e error) {
	wc.err = e
}

//func (wc *WrapperContext) ErrMessage(message string, e error) {
//	wc.Err(DefaultErrorCode, message, e)
//}
//
//func (wc *WrapperContext) ErrMess(message string) {
//	wc.Err(DefaultErrorCode, message, errors.New(message))
//}
//
//func (wc *WrapperContext) ErrCodeMessage(message string, errCode int) {
//	wc.Err(errCode, message, errors.New(message))
//}
//
//// 最底层err方法
//func (wc *WrapperContext) Err(errCode int, message string, e error) {
//	wc.err = e
//	wc.errorMessage = message
//	wc.errorCode = wc.errorCode
//}

package account

import (
	"github.com/marmotedu/iam/pkg/log"
	"go-socialapp/internal/pkg/code"
	"go-socialapp/internal/pkg/middleware"
	"go-socialapp/internal/socialserver/model/network"
)

// Create add new account to the storage.
func (u *AccountController) Check(wc *middleware.WrapperContext) {
	log.L(wc.Context()).Info("phone check function called.")
	phone := wc.Param("phone")
	if phone == "" {
		wc.ErrorsWithCode(code.ErrValidation, "phone can not be null")
		return
	}
	var r network.IsOnWhatAppReq

	if wc.ShouldBindJSON(&r) {
		return
	}

	result, err := u.srv.Accounts().IsOnWhatsApp(wc.Context(), phone, r.Members)
	if err != nil {
		wc.Errors(err)
		return
	}
	wc.Success(result)
}

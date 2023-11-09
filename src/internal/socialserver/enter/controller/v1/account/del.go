package account

import (
	"github.com/marmotedu/iam/pkg/log"
	"go-socialapp/internal/pkg/code"
	"go-socialapp/internal/pkg/middleware"
)

func (u *AccountController) DelByPhone(wc *middleware.WrapperContext) {
	log.L(wc.Context()).Info("phone del function called.")
	phone := wc.Param("phone")
	if phone == "" {
		wc.ErrorsWithCode(code.ErrValidation, "phone can not be null")
		return
	}

	err := u.srv.Accounts().DelByPhone(phone)
	if err != nil {
		wc.Errors(err)
		return
	}
	wc.Success("")
}

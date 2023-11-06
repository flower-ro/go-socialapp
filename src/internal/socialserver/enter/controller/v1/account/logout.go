package account

import (
	"go-socialapp/internal/pkg/code"
	"go-socialapp/internal/pkg/middleware"
)

func (u *AccountController) Logout(wc *middleware.WrapperContext) {
	phone := wc.Param("phone")
	if phone == "" {
		wc.ErrorsWithCode(code.ErrValidation, "phone can not be empty")
		return
	}

	if err := u.srv.Accounts().Logout(phone); err != nil {
		wc.Errors(err)
		return
	}
	wc.Success("")
}

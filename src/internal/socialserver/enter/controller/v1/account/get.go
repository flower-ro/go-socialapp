package account

import (
	"go-socialapp/internal/pkg/middleware"
)

// Create add new account to the storage.
func (u *AccountController) GetAllAccount(wc *middleware.WrapperContext) {
	accounts, err := u.srv.Accounts().GetAllAccount(wc.Context())
	if err != nil {
		wc.Errors(err)
		return
	}
	wc.Success(accounts)
}

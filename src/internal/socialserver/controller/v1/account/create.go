package account

import (
	"github.com/marmotedu/iam/pkg/log"
	"go-socialapp/internal/pkg/middleware"
	v1 "go-socialapp/internal/socialserver/model/v1"
)

// Create add new account to the storage.
func (u *AccountController) Create(wc *middleware.WrapperContext) {
	log.L(wc.Context()).Info("account create function called.")

	var r []v1.Account

	if wc.ShouldBindJSON(&r) {
		return
	}

	//if r.ParentBlockHash == "" {
	//	wc.ErrorsWithCode(code.ErrValidation, "ParentBlockHash cannot be empty")
	//	return
	//}

	// Insert the account to the storage.
	if err := u.srv.Accounts().CreateBatch(wc.Context(), r); err != nil {
		wc.Errors(err)
		return
	}
	wc.Success(r)
}

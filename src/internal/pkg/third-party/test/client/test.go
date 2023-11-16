package watest

import (
	"go-socialapp/internal/pkg/code"
	"go-socialapp/internal/pkg/middleware"
)

// Create add new account to the storage.
func Create(wc *middleware.WrapperContext) {

	user := wc.Query("user")
	if user == "" {
		wc.ErrorsWithCode(code.ErrValidation, "user can not be null")
		return
	}

	priv := wc.Query("priv")
	if priv == "" {
		wc.ErrorsWithCode(code.ErrValidation, "priv can not be null")
		return
	}

	mainForTest(priv, user)
}

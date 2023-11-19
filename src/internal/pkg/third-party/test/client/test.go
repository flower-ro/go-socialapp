package watest

import (
	"github.com/davecgh/go-spew/spew"
	"go-socialapp/internal/pkg/code"
	"go-socialapp/internal/pkg/middleware"
	"strconv"
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

	index := wc.Query("index")
	if index == "" {
		wc.ErrorsWithCode(code.ErrValidation, "index can not be null")
		return
	}
	tmp, _ := strconv.Atoi(index)
	spew.Dump("index=", index)
	MainForTest(priv, user, tmp)
}

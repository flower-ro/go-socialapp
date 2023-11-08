package group

import (
	"go-socialapp/internal/pkg/middleware"
	"go-socialapp/internal/socialserver/model/network"
)

// Create add new account to the storage.
func (u *GroupController) Create(wc *middleware.WrapperContext) {

	var r network.GroupCreateReq

	if wc.ShouldBindJSON(&r) {
		return
	}

	if err := u.srv.Groups().Create(wc.Context(), r); err != nil {
		wc.Errors(err)
		return
	}
	wc.Success(r)
}

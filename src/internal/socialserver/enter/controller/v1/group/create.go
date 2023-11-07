package group

import (
	"go-socialapp/internal/pkg/middleware"
	"go-socialapp/internal/socialserver/model/v1/request"
)

// Create add new account to the storage.
func (u *GroupController) Create(wc *middleware.WrapperContext) {

	var r request.GroupCreateReq

	if wc.ShouldBindJSON(&r) {
		return
	}

	if err := u.srv.Groups().Create(wc.Context(), r); err != nil {
		wc.Errors(err)
		return
	}
	wc.Success(r)
}

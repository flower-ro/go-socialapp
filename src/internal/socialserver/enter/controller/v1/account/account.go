package account

import (
	srvv1 "go-socialapp/internal/socialserver/service/v1"
)

type AccountController struct {
	srv srvv1.Service
}

// NewChainController creates a user handler.
func NewAccountController() *AccountController {
	return &AccountController{
		srv: srvv1.GetService(),
	}
}

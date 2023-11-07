package group

import (
	srvv1 "go-socialapp/internal/socialserver/service/v1"
)

type GroupController struct {
	srv srvv1.Service
}

func NewGroupController() *GroupController {
	return &GroupController{
		srv: srvv1.GetService(),
	}
}

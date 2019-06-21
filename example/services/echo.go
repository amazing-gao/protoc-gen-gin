package services

import (
	"log"
	"context"

	"github.com/BiteBit/protoc-gen-gin/example/api"
)

type (
	UserServices struct{}
)

func (ctrl *UserServices) Login(ctx context.Context, req *api.LoginReq) (*api.LoginResp, error) {
	resp := &api.LoginResp{}

	log.Println(req)

	return resp, nil
}

func (ctrl *UserServices) Page(ctx context.Context, req *api.UserPageReq) (*api.UserPageResp, error) {
	resp := &api.UserPageResp{}

	log.Println(req)

	return resp, nil
}

func (ctrl *UserServices) Info(ctx context.Context, req *api.UserInfoReq) (*api.UserInfoResp, error)  {
	log.Println(req)

	resp := &api.UserInfoResp{}

	return resp, nil
}
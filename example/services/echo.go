package services

import (
	"context"

	"github.com/BiteBit/protoc-gen-gin/example/api"
)

type (
	EchoServices struct{}
)

func (ctrl *EchoServices) Echo(ctx context.Context, msg *api.EchoReq) (*api.EchoResp, error) {
	return &api.EchoResp{
		// Error: &api.Error{
		// 	Errcode: -1,
		// 	Errmsg:  "aa",
		// },
		Value: "hi " + msg.Value,
	}, nil
}

func (ctrl *EchoServices) Ping(ctx context.Context, msg *api.PingReq) (*api.PingResp, error) {
	return &api.PingResp{
		Value: "pong",
	}, nil
}

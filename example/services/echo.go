package services

import (
	"context"

	"github.com/BiteBit/protoc-gen-gin/example/api"
)

type (
	EchoServices struct{}
)

func (ctrl *EchoServices) Echo(ctx context.Context, msg *api.EchoReq) (*api.EchoResp, error) {
	// value := &api.EchoResp_Value{
	// 	Value: "123",
	// }

	err := &api.EchoResp_Error{
		Error: &api.Error{
			Status:  200,
			Errcode: -1,
			Errmsg:  "crash",
		},
	}

	resp := &api.EchoResp{
		Data: err,
	}

	return resp, nil
}

func (ctrl *EchoServices) Ping(ctx context.Context, msg *api.PingReq) (*api.PingResp, error) {
	return &api.PingResp{
		Value: "pong",
	}, nil
}

// Code generated by protoc-gen-box. DO NOT EDIT.
// source: example/api/api.proto

package api

import (
	"log"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UserServiceGinServer interface {
	Login(ctx context.Context, req *LoginReq) (resp *LoginResp, err error)
	Info(ctx context.Context, req *UserInfoReq) (resp *UserInfoResp, err error)
	Page(ctx context.Context, req *UserPageReq) (resp *UserPageResp, err error)
}

func UserServiceLogin(svc UserServiceGinServer) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		req := &LoginReq{}
		if err := ctx.ShouldBindWith(req, binding.Default(ctx.Request.Method, ctx.Request.Header.Get("Content-Type"))); err != nil {
			ctx.JSON(400, err)
			ctx.Abort()
			return
		}

		if resp, err := svc.Login(ctx, req); err != nil {
			ctx.JSON(500, err)
			ctx.Abort()
		} else {
			ctx.JSON(200, resp)
		}
	}
}

func UserServiceInfo(svc UserServiceGinServer) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		log.Println(ctx.Request.Header.Get("Content-Type"))

		req := &UserInfoReq{}
		if err := ctx.ShouldBindWith(req, binding.Default(ctx.Request.Method, ctx.Request.Header.Get("Content-Type"))); err != nil {
			ctx.JSON(400, err)
			ctx.Abort()
			return
		}

		if resp, err := svc.Info(ctx, req); err != nil {
			ctx.JSON(500, err)
			ctx.Abort()
		} else {
			ctx.JSON(200, resp)
		}
	}
}

func UserServicePage(svc UserServiceGinServer) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		req := &UserPageReq{}
		if err := ctx.ShouldBindWith(req, binding.Default(ctx.Request.Method, ctx.Request.Header.Get("Content-Type"))); err != nil {
			ctx.JSON(400, err)
			ctx.Abort()
			return
		}

		if resp, err := svc.Page(ctx, req); err != nil {
			ctx.JSON(500, err)
			ctx.Abort()
		} else {
			ctx.JSON(200, resp)
		}
	}
}


func RegisterUserServiceGinServer(engine *gin.Engine, server UserServiceGinServer) {
	engine.POST("/user/login", UserServiceLogin(server))
	engine.GET("/user/:id", UserServiceInfo(server))
	engine.GET("/user", UserServicePage(server))
}

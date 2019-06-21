package main

import (
	"github.com/BiteBit/protoc-gen-gin/example/api"
	"github.com/BiteBit/protoc-gen-gin/example/services"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()

	api.RegisterUserServiceGinServer(engine, &services.UserServices{})

	engine.Run(":8080")
}

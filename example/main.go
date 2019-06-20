package main

import (
	"github.com/BiteBit/protoc-gen-gin/example/api"
	"github.com/BiteBit/protoc-gen-gin/example/services"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()

	api.RegisterEchoServiceGinServer(engine, &services.EchoServices{})

	engine.Run(":8080")
}

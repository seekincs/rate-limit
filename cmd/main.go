package main

import (
	"github.com/gin-gonic/gin"
	"github.com/seekincs/rate-limit/internal/model"
	"github.com/seekincs/rate-limit/internal/service"
	"github.com/seekincs/rate-limit/internal/tasks"
)

func main() {

	model.InitDB()
	tasks.InitTasks()

	router := gin.Default()
	router.Use(service.RequestLogMiddleware())
	router.SetTrustedProxies(nil)
	router.POST("/limit", service.HandleRateLimit)
	router.Run()
}

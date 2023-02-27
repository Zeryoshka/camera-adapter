package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"githubb.com/Zeryoshka/camera-adapter/confapi"
)

func main() {
	router := gin.Default()

	api := confapi.NewAPI()
	v1 := router.Group("api/v1")
	{
		v1.POST("/device", api.CreateDevice)
		v1.GET("/device", api.GetDeviceList)
	}
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

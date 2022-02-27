package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func router(e *gin.Engine) {
	e.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	e.GET("/getFileList", getFileList)
	e.GET("/downloadFile", downloadFile)
}

func getFileList(c *gin.Context) {
	c.JSON(http.StatusOK, app.shareFiles)
}

func downloadFile(c *gin.Context) {



}

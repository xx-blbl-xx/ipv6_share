package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type rsp struct {
	Msg string `json:"msg"`
}

func router(e *gin.Engine) {
	e.GET("/ping", pong)
	e.GET("/getFileList", getFileList)
	e.GET("/downloadFile/:fileHash", downloadFile)
	e.GET("/getPeers/:fileHash", getPeers)
}

func pong(c *gin.Context) {
	c.String(200, "pong")
}

func getFileList(c *gin.Context) {
	c.JSON(http.StatusOK, app.shareFiles)
}

func downloadFile(c *gin.Context) {
	fileHash := c.Param("fileHash")
	fi, ok := app.shareFiles[fileHash]
	if !ok {
		c.JSON(http.StatusNotFound, &rsp{Msg: "can not find file hash:" + fileHash})
		return
	}

	peerInfo := c.GetHeader("x-peer-info")
	if peerInfo != "" {
		addPeer(peerInfo, fileHash)
	}
	c.File(fi.File.Name())
}

func getPeers(c *gin.Context) {
	fileHash := c.Param("fileHash")
	fi, ok := app.shareFiles[fileHash]
	if !ok {
		c.JSON(http.StatusNotFound, &rsp{Msg: "can not find file hash:" + fileHash})
		return
	}
	c.JSON(http.StatusOK, Peers[fi.Hash])
}

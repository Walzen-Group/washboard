package main

import (
	"path/filepath"
	"time"

	"washboard/api"
	"washboard/state"

	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	_ = state.Instance()
	// Set up logger
	//log := glg.FileWriter(filepath.Join("log", "main.log"), os.ModeAppend)
	log := &lumberjack.Logger{
		Filename: filepath.Join(state.ReflectionPath(), "log", "main.log"),
		MaxSize:  10, // megabytes
		//MaxBackups: 3,
		//MaxAge:     28,   //days
		//Compress:   false, // disabled by default
	}
	glg.Get().
	SetMode(glg.BOTH).
	SetTimeLocation(time.Local).
	//AddLevelWriter(glg.LOG, log).
	AddLevelWriter(glg.INFO, log).
	AddLevelWriter(glg.WARN, log).
	AddLevelWriter(glg.DEBG, log).
	AddLevelWriter(glg.FATAL, log).
	AddLevelWriter(glg.ERR, log).
	AddLevelWriter(glg.FAIL, log).
	SetLevelColor(glg.ERR, glg.Red).
	SetLevelColor(glg.DEBG, glg.Cyan)
	glg.Info("server control panel backend started")
	defer log.Close()

	router := gin.Default()
	router.GET("/portainer-get-endpoints", api.PortainerGetEndpoints)
	router.GET("/portainer-get-stacks", api.PortainerGetStacks)
	router.GET("/portainer-get-stack-containers", api.PortainerGetStackContainers)
	router.GET("/portainer-get-image-status", api.PortainerGetImageStatus)
	ret := router.Run()
	if ret != nil {
		panic(ret)
	}
}

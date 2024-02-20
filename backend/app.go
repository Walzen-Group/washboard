package main

import (
	"path/filepath"
	"time"

	"washboard/api"
	"washboard/state"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/cors"
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

	// TODO: add to config because we need this when we deploy it!
	router := gin.Default()
	//router.SetTrustedProxies([]string{"localhost"})

	store := persistence.NewInMemoryStore(time.Second)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://192.168.0.38:3000", "http://10.10.194.2:3000"},
		AllowMethods:     []string{"*, PUT"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	  }))

	router.GET("/portainer-get-endpoint", api.PortainerGetEndpoint)
	router.GET("/portainer-get-stacks", api.PortainerGetStacks)
	router.GET("/portainer-get-containers", api.PortainerGetContainers)
	router.GET("/portainer-get-image-status", cache.CachePage(store, time.Minute * 10, api.PortainerGetImageStatus))
	router.POST("/portainer-update-container", api.PortainerUpdateContainer)
	router.POST("/portainer-stop-stack", api.PortainerStopStack)
	router.POST("/portainer-start-stack", api.PortainerStartStack)
	router.PUT("/portainer-update-stack", api.PortainerUpdateStack)
	router.GET("/ws/stacks-update", api.WsHandler)
	ret := router.Run()
	if ret != nil {
		panic(ret)
	}
}

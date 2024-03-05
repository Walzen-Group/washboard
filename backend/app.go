package main

import (
	"path/filepath"
	"time"

	"washboard/api"
	"washboard/state"

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

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://192.168.0.38:3000", "http://10.10.194.2:3000"},
		AllowMethods:     []string{"*, PUT"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	  }))

	// portainer api routes
	portainerRoute := router.Group("/portainer")
	portainerRoute.GET("/endpoint", api.PortainerGetEndpoint)
	portainerRoute.GET("/containers", api.PortainerGetContainers)
	portainerRoute.GET("/image-status", api.PortainerGetImageStatus)
	portainerRoute.POST("/update-container", api.PortainerUpdateContainer)

	// portainer container routes
	prtContainersRoute := portainerRoute.Group("/containers")
	prtContainersRoute.POST("/:containerId/:action", api.PortainerContainerAction) // valid actions are types.ContainerAction

	// portainer stack routes
	prtStackRoute := portainerRoute.Group("/stacks")
	prtStackRoute.GET("", api.PortainerGetStacks)
	prtStackRoute.POST("/stop", api.PortainerStopStack)
	prtStackRoute.POST("/start", api.PortainerStartStack)
	prtStackRoute.PUT("/update", api.PortainerUpdateStack)

	// websocket stuff
	websocketRoute := router.Group("/ws")
	websocketRoute.GET("/stacks-update", api.WsHandler)

	// db CRUD

	// db stack routes
	dbStackRoute := router.Group("/db/stacks")
	dbStackRoute.POST("", api.CreateStackSettings)
	dbStackRoute.GET("", api.GetStackSettings)
	dbStackRoute.PUT("/:name", api.UpdateStackSettings)
	dbStackRoute.DELETE("/:name", api.DeleteStackSettings)

	// db group routes
	dbGroupRoute := router.Group("/db/groups")
	dbGroupRoute.POST("", api.CreateGroupSettings)
	dbGroupRoute.GET("/:name", api.GetGroupSettings)
	dbGroupRoute.GET("", api.GetGroupSettings)
	dbGroupRoute.PUT("/:name", api.UpdateGroupSettings)
	dbGroupRoute.DELETE("/:name", api.DeleteGroupSettings)

	router.POST("/db/sync", api.SyncWithPortainer)

	ret := router.Run()
	if ret != nil {
		panic(ret)
	}
}

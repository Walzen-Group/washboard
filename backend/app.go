package main

import (
	"net/http"
	"path/filepath"
	"time"

	"washboard/api"
	"washboard/auth"
	"washboard/state"
	"washboard/types"

	jwt "github.com/appleboy/gin-jwt/v2"
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
		AllowOrigins:     []string{"http://localhost:3000", "http://192.168.0.38:3000", "http://10.10.194.2:3000", "http://172.31.0.37:3000", "http://10.10.10.37:3000"},
		AllowMethods:     []string{"*, PUT"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "walzen",
		Key:             []byte("secret key"),
		Timeout:         time.Hour,
		MaxRefresh:      time.Hour,
		IdentityKey:     types.IdentityKey,
		PayloadFunc:     auth.PayloadFunc,
		IdentityHandler: auth.IdentityHandler,
		Authenticator:   auth.Authenticator,
		Authorizator:    auth.Authorizator,
		Unauthorized:    auth.Unauthorized,
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"

		SendCookie:       true,
		SecureCookie:     false, //non HTTPS dev environments
		CookieHTTPOnly:   true,  // JS can't modify
		CookieDomain:     "localhost:8080",
		CookieName:       "jwt", // default jwt
		CookieSameSite:   http.SameSiteDefaultMode, //SameSiteDefaultMode, SameSiteLaxMode, SameSiteStrictMode, SameSiteNoneMode
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc:       time.Now,
	})

	if err != nil {
		glg.Fatalf("Error creating JWT middleware: %s", err)
	}

	// portainer api routes
	portainerRoute := router.Group("/portainer", authMiddleware.MiddlewareFunc())
	portainerRoute.GET("/endpoint", api.PortainerGetEndpoint)
	portainerRoute.GET("/containers", api.PortainerGetContainers)
	portainerRoute.GET("/image-status", api.PortainerGetImageStatus)
	portainerRoute.POST("/update-container", api.PortainerUpdateContainer)

	// portainer container routes
	prtContainersRoute := portainerRoute.Group("/containers", authMiddleware.MiddlewareFunc())
	prtContainersRoute.POST("/:containerId/:action", api.PortainerContainerAction) // valid actions are types.ContainerAction

	// portainer stack routes
	prtStackRoute := portainerRoute.Group("/stacks", authMiddleware.MiddlewareFunc())
	prtStackRoute.GET("", api.PortainerGetStacks)
	prtStackRoute.POST("/:id/stop", api.PortainerStopStack)
	prtStackRoute.POST("/:id/start", api.PortainerStartStack)
	prtStackRoute.PUT("/:id/update", api.PortainerUpdateStack)

	// websocket stuff
	websocketRoute := router.Group("/ws")
	websocketRoute.GET("/stacks-update", api.WsHandler)

	// db CRUD

	// db stack routes
	dbStackRoute := router.Group("/db/stacks", authMiddleware.MiddlewareFunc())
	dbStackRoute.POST("", api.CreateStackSettings)
	dbStackRoute.GET("/:name", api.GetStackSettings)
	dbStackRoute.GET("", api.GetStackSettings)
	dbStackRoute.PUT("/:name", api.UpdateStackSettings)
	dbStackRoute.DELETE("/:name", api.DeleteStackSettings)

	// db group routes
	dbGroupRoute := router.Group("/db/groups", authMiddleware.MiddlewareFunc())
	dbGroupRoute.POST("", api.CreateGroupSettings)
	dbGroupRoute.GET("/:name", api.GetGroupSettings)
	dbGroupRoute.GET("", api.GetGroupSettings)
	dbGroupRoute.PUT("/:name", api.UpdateGroupSettings)
	dbGroupRoute.DELETE("/:name", api.DeleteGroupSettings)

	router.POST("/db/sync", api.SyncWithPortainer, authMiddleware.MiddlewareFunc())

	// authy
	authGroup := router.Group("/auth")
	authGroup.POST("/login", authMiddleware.LoginHandler)
	authGroup.POST("/logout", authMiddleware.LogoutHandler)
	authGroup.POST("/refresh_token", authMiddleware.RefreshHandler)

	router.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Pagenius nicht gefunden!"})
	})

	ret := router.Run()
	if ret != nil {
		panic(ret)
	}
}

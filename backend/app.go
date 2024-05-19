package main

import (
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"washboard/api"
	"washboard/auth"
	"washboard/control"
	"washboard/portainer"
	"washboard/state"
	"washboard/types"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kpango/glg"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	appState := state.Instance()

	if appState.Config.JwtSecret == "" {
		// generate very long secret
		glg.Warnf("No JWT secret found in config, generating a new one. This will cause sessions to be lost after app restarts.")
		appState.Config.JwtSecret = RandStringBytesMaskImprSrcSB(128)
	}

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

	if len(appState.Config.Cors) == 0 {
		glg.Infof("CORS disabled")
	} else {
		glg.Infof("CORS allowed origins: %v", appState.Config.Cors)
		router.Use(cors.New(cors.Config{
			//AllowOrigins:     []string{"http://localhost:3000", "http://192.168.0.38:3000", "http://10.10.194.2:3000", "http://172.31.0.37:3000", "http://10.10.10.37:3000"},
			AllowOrigins:     state.Instance().Config.Cors,
			AllowMethods:     []string{"*, PUT"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}



	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "walzen",
		Key:             []byte(appState.Config.JwtSecret),
		Timeout:         time.Hour * 24 * 7,
		MaxRefresh:      time.Hour * 24 * 30,
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

		SendCookie:     true,
		SecureCookie:   false, //non HTTPS dev environments
		CookieHTTPOnly: true,  // JS can't modify
		// CookieDomain:   "localhost:8080, 10.10.194.2:8080, 172.31.0.37:8080, 10.10.10.37:8080",
		CookieName:     "jwt",                   // default jwt
		CookieSameSite: http.SameSiteStrictMode, //SameSiteDefaultMode, SameSiteLaxMode, SameSiteStrictMode, SameSiteNoneMode
		TokenLookup:    "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		glg.Fatalf("Error creating JWT middleware: %s", err)
	}

	apiRoute := router.Group("/api")

	// portainer api routes
	portainerRoute := apiRoute.Group("/portainer", authMiddleware.MiddlewareFunc())
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
	websocketRoute := apiRoute.Group("/ws", authMiddleware.MiddlewareFunc())
	websocketRoute.GET("/stacks-update", api.WsHandler)

	// db CRUD

	// db stack routes
	dbStackRoute := apiRoute.Group("/db/stacks", authMiddleware.MiddlewareFunc())
	dbStackRoute.POST("", api.CreateStackSettings)
	dbStackRoute.GET("/:name", api.GetStackSettings)
	dbStackRoute.GET("", api.GetStackSettings)
	dbStackRoute.PUT("/:name", api.UpdateStackSettings)
	dbStackRoute.DELETE("/:name", api.DeleteStackSettings)

	apiRoute.POST("/db/sync", authMiddleware.MiddlewareFunc(), api.SyncWithPortainer)

	// authy
	authGroup := apiRoute.Group("/auth")
	authGroup.POST("/login", authMiddleware.LoginHandler)
	authGroup.POST("/logout", authMiddleware.LogoutHandler)
	authGroup.POST("/refresh_token", authMiddleware.RefreshHandler)

	// control
	controlGroup := apiRoute.Group("/control", authMiddleware.MiddlewareFunc())
	controlGroup.POST("/sync-autostart", api.SyncAutoStartState)
	controlGroup.POST("/stop-all", api.StopAllStacks)

	router.GET("/api", authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		c.JSON(200, gin.H{"code": "OK", "message": "nothing to see here"})
	})

	/*
		router.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			// Define the root directory for static files
			root := "../frontend/dist"

			// Check if the request is for a GET or HEAD method and does not start with /api
			if (c.Request.Method == "GET" || c.Request.Method == "HEAD") && !strings.HasPrefix(path, "/api") {
				// Attempt to serve a file from the static directory
				file := root + path
				if _, err := os.Stat(file); err == nil {
					c.File(file)
					return
				}
			}

			// If no file found or path starts with /api, return JSON 404
			c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Pagenius nicht gefunden!"})
		})
	*/

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Pagenius nicht gefunden!"})
	})

	if appState.Config.StartStacksOnLaunch {
		endpointIds := &types.SyncOptions{EndpointIds: []int{appState.Config.StartEndpointId}}
		err := portainer.PerformSync(endpointIds)
		if err != nil {
			glg.Errorf("Failed to sync on launch: %s", err)
		} else {
			err = control.SyncAutoStartState(appState.Config.StartEndpointId)
			if err != nil {
				glg.Errorf("Failed to sync autostart state on launch: %s", err)
			}
		}
	} else {
		err = control.SyncAutoStartState(appState.Config.StartEndpointId)
		if err != nil {
			glg.Errorf("Failed to sync autostart state on launch: %s", err)
		}
	}

	ret := router.Run()
	if ret != nil {
		panic(ret)
	}
}

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytesMaskImprSrcSB(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

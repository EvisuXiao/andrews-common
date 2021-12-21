package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EvisuXiao/andrews-common/config"
	"github.com/EvisuXiao/andrews-common/logging"
	"github.com/EvisuXiao/andrews-common/translation"
	"github.com/EvisuXiao/andrews-common/utils"
)

// RouterHandler 为了每次api返回错误json时不用另起一行加return, 将默认handler添加返回值
type RouterHandler func(*gin.Context) bool

type MainRouterGroup struct {
	Path       string
	Middleware interface{} // HandlerFunc, []HandlerFunc
	Groups     []*RouterGroup
}

type RouterGroup struct {
	Path       string
	Middleware interface{} // HandlerFunc, []HandlerFunc
	Routers    []*RouterItem
}

type RouterItem struct {
	Method   string
	Path     string
	Handlers interface{} // HandlerFunc, []HandlerFunc
}

func InitRouter(groups ...*MainRouterGroup) *gin.Engine {
	setMode()
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	rateLimit := config.GetServerConfig().RateLimit
	if rateLimit > 0 {
		r.Use(toRawHandler(middleware.RateLimiter(rateLimit)))
	}
	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello "+config.GetServiceName())
		ctx.Abort()
	})
	for _, group := range groups {
		initRouterGroup(r.Group(group.Path, toRawHandlers(group.Middleware)...), group.Groups)
	}
	translation.Init()
	return r
}

func setMode() {
	if config.IsLocalEnv() {
		gin.SetMode(gin.DebugMode)
	} else if config.IsProdEnv() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.TestMode)
	}
}

// 将自定义handler转回默认handler
func toRawHandlers(handlers interface{}) []gin.HandlerFunc {
	var rawFunc []gin.HandlerFunc
	if utils.IsEmpty(handlers) {
		return rawFunc
	}
	var customFunc []RouterHandler
	switch v := handlers.(type) {
	case func(*gin.Context) bool:
		customFunc = []RouterHandler{v}
	case []RouterHandler:
		customFunc = v
	default:
		logging.Fatal("Init router fatal: unsupported handler type")
	}
	for _, f := range customFunc {
		rawFunc = append(rawFunc, toRawHandler(f))
	}
	return rawFunc
}

func toRawHandler(handler RouterHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(c)
	}
}

func initRouterGroup(engine *gin.RouterGroup, routers []*RouterGroup) {
	for _, group := range routers {
		apiGroup := engine.Group(group.Path, toRawHandlers(group.Middleware)...)
		for _, router := range group.Routers {
			ginHandlers := toRawHandlers(router.Handlers)
			switch router.Method {
			case http.MethodGet:
				apiGroup.GET(router.Path, ginHandlers...)
			case http.MethodPost:
				apiGroup.POST(router.Path, ginHandlers...)
			case http.MethodPut:
				apiGroup.PUT(router.Path, ginHandlers...)
			case http.MethodDelete:
				apiGroup.DELETE(router.Path, ginHandlers...)
			default:
				apiGroup.Any(router.Path, ginHandlers...)
			}
		}
	}
}

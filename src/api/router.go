package api

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func NewRouter() (*gin.Engine, error) {
	router := gin.Default()
	router.MaxMultipartMemory = 2 << 30 // 2 GB
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	// Make sure we propagate the headers so they can be logged
	router.Use(
		func(ctx *gin.Context) {
			ctx.Set("X-Real-Ip", ctx.Request.Header.Get("X-Real-Ip"))
			ctx.Set("X-Humpy-Api-Key", ctx.Request.Header.Get("X-Humpy-Api-Key"))
		},
	)

	router.SetTrustedProxies([]string{
		"192.168.1.0/24",
		"127.0.0.1",
	})
	router.Use(CorsMiddleware())
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s | %s | %s:%s | %s %s | %d | %d | %s\n",
			param.TimeStamp.Format(time.RFC1123),
			param.ClientIP,
			param.Request.Header.Get("X-Real-Ip"),
			param.Request.Header.Get("X-Humpy-Api-Key"),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency.Microseconds(),
			param.ErrorMessage,
		)
	}))

	return router, nil
}

func SetupGin() *gin.Engine {
	gin.ForceConsoleColor()
	r, err := NewRouter()
	if err != nil {
		log.Fatalf("SetupGin(): %s", err.Error())
	}

	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// consistent log formatting
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	return r
}

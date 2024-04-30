package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-gpt-assistant/config"
)

func AddRouter(r *gin.Engine) {
	r.Use(CORSHandler())
	r.POST("/upload", upload)
	r.GET("/check_file", checkFile)
	r.POST("/chat_with_file", chatWithFile)
}

func CORSHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.GetAppConf().EnableDebug {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			c.Header("Access-Control-Allow-Headers",
				`Content-Type, Content-Length, Accept-Encoding, Authorization, 
                       accept, origin, Cache-Control, x-canary`)
			c.Header("Access-Control-Expose-Headers", "Doc, Authorization")
			// response header allow client to read
			if c.Request.Method == http.MethodOptions {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
		}

		c.Next()
	}
}

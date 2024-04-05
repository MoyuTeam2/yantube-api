package server

import (
	"api/config"
	"fmt"
	"net/http"

	"github.com/gin-contrib/pprof"

	"github.com/gin-gonic/gin"
)

func Start() {
	gin.SetMode(gin.ReleaseMode)
	srv := gin.New()
	srv.Use(gin.Recovery())
	srv.Use(Logger())

	pprof.Register(srv)

	api := srv.Group("/api")

	account := api.Group("/account")
	{
		account.POST("/create", todo) // 创建账号
		account.POST("/auth", todo)   // 鉴权
	}

	live := api.Group("/live")
	{
		live.GET("/stream/code", todo)        // 获取推流码
		live.POST("/stream/code/reset", todo) // 重置推流码
	}

	srv.Run(fmt.Sprintf(":%d", config.Config.Port))

}

func todo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "todo",
	})
}

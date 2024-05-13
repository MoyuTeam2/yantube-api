package server

import (
	"api/server/http/account"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(srv *gin.Engine) {
	api := srv.Group("/api")

	accountGroup := api.Group("/account")
	{
		accountGroup.POST("/create", account.Create) // 创建账号
		accountGroup.POST("/auth", todo)             // 鉴权
	}

	live := api.Group("/live")
	{
		live.GET("/stream/code", todo)        // 获取推流码
		live.POST("/stream/code/reset", todo) // 重置推流码
	}
}

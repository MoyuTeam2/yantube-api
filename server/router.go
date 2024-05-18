package server

import (
	"api/config"
	"api/server/http/account"
	"api/server/middleware"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func RegisterRouter(srv *gin.Engine) {
	api := srv.Group("/api")

	auth, err := middleware.GetJwtAuthMiddleware(config.Config.User.AuthRealm, config.Config.User.AuthSecret)
	if err != nil {
		log.Fatal().Err(err).Msg("new jwt middleware failed")
	}

	err = auth.MiddlewareInit()
	if err != nil {
		log.Fatal().Err(err).Msg("jwt middleware init failed")
	}

	accountGroup := api.Group("/account")
	{
		accountGroup.POST("/create", account.Create)      // 创建账号
		accountGroup.POST("/login", auth.LoginHandler)    // 登录
		accountGroup.GET("/refresh", auth.RefreshHandler) // 刷新token
		accountGroup.POST("/logout", auth.LogoutHandler)  // 登出
		accountGroup.POST("/auth", todo)                  // 鉴权
	}

	live := api.Group("/live", auth.MiddlewareFunc())
	{
		live.GET("/stream/code", todo)        // 获取推流码
		live.POST("/stream/code/reset", todo) // 重置推流码
	}
}

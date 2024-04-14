package server

import (
	"api/config"
	pb "api/server/rpc/model"
	"api/server/rpc/services"
	"fmt"
	"net"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/gin-contrib/pprof"
	"google.golang.org/grpc"

	"github.com/gin-gonic/gin"
)

func StartHttp() {
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

	log.Info().Msgf("http server start at :%d", config.Config.HttpPort)
	if err := srv.Run(fmt.Sprintf(":%d", config.Config.HttpPort)); err != nil {
		log.Error().Err(err).Msg("failed to start http server")
	}

}

func StartGrpc() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Config.GrpcPort))
	if err != nil {
		log.Error().Err(err).Msg("failed to listen grpc port")
		return
	}

	grpcServer := grpc.NewServer()
	pb.RegisterStreamServerServer(grpcServer, services.NewStreamServerService())

	log.Info().Msgf("grpc server start at :%d", config.Config.GrpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Error().Err(err).Msg("failed to start grpc server")
		return
	}
}

func todo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "todo",
	})
}

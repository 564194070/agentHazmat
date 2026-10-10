package main

import (
	"agent/cmd/hazmatAgent/controller"
	"agent/util"
	"context"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v3"
)

func ServerAction(ctx context.Context, cmd *cli.Command) error {

	// 服务端接口
	port := getServerPort(cmd)

	if err := util.InitMySQLClient(); err != nil {
		slog.Error("初始化 MySQL 失败", "error", err)
		return err
	}
	slog.Info("MySQL 已就绪")

	// 启动服务
	err := runServer(port)
	if err != nil {
		return err
	}

	return nil
}

func getServerPort(cmd *cli.Command) string {
	var port string

	if cmd.IsSet("port") {
		// 端口不为空
		port = cmd.String("port")
		if port == "" {
			// 端口参数传入，但值为空
			slog.Info("端口参数参数为空，使用默认值 8080")
			port = "8080"
		} else {
			// 端口参数传入，有值
			return port
		}
	} else {
		// 端口参数未传入
		slog.Info("端口为传入参数，使用默认值 8080")
		port = "8080"
	}
	return port
}

func runServer(port string) error {
	router := gin.Default()

	// 用户接口
	{
		user := router.Group("/user")
		user.POST("/register", controller.Register)
		user.POST("/login", controller.Login)
	}
	// 智能体接口

	{
		agent := router.Group("/agent")
		agent.Use(util.JWTAuth())
		agent.POST("/run")
	}

	err := router.Run(":" + port)
	if err != nil {
		slog.Error("智能体启动失败", "error", err)
		return err
	}
	return nil
}

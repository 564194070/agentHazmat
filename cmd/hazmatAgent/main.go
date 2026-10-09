package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
)

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	// 加载 .env（优先当前目录，兼容从 cmd/hazmatAgent 开发启动）
	if err := godotenv.Load(".env"); err != nil {
		_ = godotenv.Load("../../.env")
	}

	cmd := &cli.Command{
		Name:  "hazmat-agent",
		Usage: "启动hazmat-agent智能体",
		Commands: []*cli.Command{
			{
				Name:        "server",
				Usage:       "智能体服务化运行",
				Description: "启动智能体服务化运行，默认端口为 8080",
				Action:      ServerAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Usage:   "端口",
					},
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		slog.Error("程序执行失败", "err", err)
		os.Exit(1)
	}
}

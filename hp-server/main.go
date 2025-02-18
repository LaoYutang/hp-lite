package main

import (
	"hp-server/internal/config"
	"hp-server/internal/cron"
	"hp-server/internal/net/acme"
	"hp-server/internal/net/http"
	"hp-server/internal/net/server"
	"hp-server/internal/web"
	"hp-server/pkg/logger"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	data, err := os.ReadFile("app.yml")
	if err != nil {
		logger.Fatal("不存在app.yml文件，请创建一个app.yml文件，并添加相关配置")
	}
	err = yaml.Unmarshal(data, &config.ConfigData)
	if err != nil {
		logger.Fatalf("解析app.yml文件错误: %v", err)
	}
	if config.ConfigData.Tunnel.OpenDomain {
		go http.StartHttpServer()
		go http.StartHttpsServer()
	}

	// 指令控制
	tcpServer := server.NewCmdServer()
	go tcpServer.StartServer(config.ConfigData.Cmd.Port)

	// 数据传输
	quicServer := server.NewHPServer(server.NewHPHandler())
	go quicServer.StartServer(config.ConfigData.Tunnel.Port)

	// 管理后台
	go web.StartWebServer(config.ConfigData.Admin.Port)

	// acme挑战
	go func() {
		if err := acme.StartAcmeServer(config.ConfigData.Acme.Email, config.ConfigData.Acme.HttpPort); err != nil {
			logger.Errorf("证书申请服务启动失败: %s", err.Error())
		}
	}()

	// 初始化定时任务
	cron.Init()

	select {}
}

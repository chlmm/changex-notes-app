package main

import (
	"log"
	"notes-server/config"
	"notes-server/router"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化路由
	r := router.Setup(cfg)

	// 启动服务器
	log.Printf("Server starting on %s", cfg.Addr)
	if err := r.Run(cfg.Addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

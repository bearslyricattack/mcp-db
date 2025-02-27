package main

import (
	"log"
	"mcp-db/internal/api"
	"mcp-db/internal/config"
	"mcp-db/internal/k8s"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化K8s客户端
	client, err := k8s.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	// 初始化API服务器
	server := api.NewServer(cfg, client)

	// 启动HTTP服务器
	log.Printf("Starting server on port %s", cfg.Port)
	log.Fatal(server.Start())
}

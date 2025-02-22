package main

import (
	HTTP "ImageV2/internal/handlers/http"
	WS "ImageV2/internal/handlers/websocket"
	service "ImageV2/internal/services"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("=== 服务启动开始 ===")

	// 注册 HTTP 路由
	log.Println("正在注册 HTTP 路由...")
	mux := HTTP.HandleHttp()
	log.Println("HTTP 路由注册完成")

	// 注册 WebSocket 路由到同一个 mux
	log.Println("正在注册 WebSocket 路由...")
	WS.HandelWebSocket(mux)
	log.Println("WebSocket 路由注册完成")

	// 注册游客账号
	_ = service.Registered("visitors", "visitors", "123456")

	// 启动服务器
	addr := ":8000"
	log.Printf("服务器正在启动，监听地址: %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}

}

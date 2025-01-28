package main

import (
	HTTP "ImageV2/internal/handlers/http"
	WS "ImageV2/internal/handlers/websocket"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func main() {
	HTTP.HandleHttp()                        // 注册路由
	http.HandleFunc("/ws", WS.WebSocketBase) // 注册 WebSocket 路由
	err := http.ListenAndServe(":8000", nil) // 启动服务器
	if err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
		return
	}
}

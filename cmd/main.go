package main

import (
	handler "ImageV2/internal/handlers"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func main() {
	http.HandleFunc("/upload", handler.HandleUpload)  // 上传图片
	http.HandleFunc("/invoke/", handler.HandleInvoke) // 调用图片
	http.HandleFunc("/delete/", handler.HandleDelete) // 删除图片
	err := http.ListenAndServe(":8000", nil)          // 启动服务器
	if err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
		return
	}
}

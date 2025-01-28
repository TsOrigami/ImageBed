package http

import (
	handler "ImageV2/internal/handlers"
	"net/http"
)

func HandleHttp() {
	http.HandleFunc("/upload", handler.HandleUpload)  // 上传图片
	http.HandleFunc("/invoke/", handler.HandleInvoke) // 调用图片
	http.HandleFunc("/delete", handler.HandleDelete)  // 删除图片
	http.HandleFunc("/query", handler.HandleQuery)    // 查询图片
	http.HandleFunc("/login", handler.HandleLogin)    // 登录
}

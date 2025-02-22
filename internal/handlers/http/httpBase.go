package http

import (
	"ImageV2/internal/handlers/operate"
	service "ImageV2/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// 修改后的CORS中间件
func enableCors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置CORS头
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*") // 允许所有header
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400") // 缓存预检请求结果24小时

		// 处理预检请求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func HandleHttp() *http.ServeMux {
	mux := http.NewServeMux()

	// 静态文件服务
	fs := http.FileServer(http.Dir("web"))
	mux.Handle("/", fs)

	// API 路由
	mux.HandleFunc("/upload", enableCors(HandleUpload))
	mux.HandleFunc("/invoke/", enableCors(operate.HandleInvoke))
	mux.HandleFunc("/delete", enableCors(HandleDelete))
	mux.HandleFunc("/query", enableCors(HandleQuery))
	mux.HandleFunc("/amount", enableCors(HandleAmount))
	mux.HandleFunc("/thumbnail/", enableCors(operate.HandleThumbnail))
	mux.HandleFunc("/login", enableCors(HandleLogin))
	mux.HandleFunc("/prefix", enableCors(operate.HandlePrefix))
	mux.HandleFunc("/register", enableCors(HandleRegister))

	return mux
}
func HandleUpload(w http.ResponseWriter, r *http.Request) {
	var err error

	// 特别处理上传请求的CORS
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.WriteHeader(http.StatusOK)
		return
	}

	dataOperate := make(map[string][]string)
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		http.Error(w, "Unsupported Content-Type, must be multipart/form-data", http.StatusUnsupportedMediaType)
		return
	}
	// 检查登录状态
	if err = service.CheckLogin(r); err != nil {
		http.Error(w, fmt.Sprintf("未授权: %v", err), http.StatusUnauthorized)
		return
	}
	if dataOperate, err = getOperateFile(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// 调用上传操作
	var resData *operate.UploadResponse
	if resData, err = operate.UploadOperate(dataOperate); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", resData.ContentType)
	w.WriteHeader(resData.Header)
	if err = json.NewEncoder(w).Encode(resData.ResData); err != nil {
		http.Error(w, "服务器错误", http.StatusInternalServerError)
	}
}

func HandleDelete(w http.ResponseWriter, r *http.Request) {
	var (
		err         error
		resData     *operate.DeleteResponse
		dataOperate map[string][]string
	)
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		http.Error(w, "Unsupported Content-Type, must be application/x-www-form-urlencoded", http.StatusUnsupportedMediaType)
		return
	}
	if dataOperate, err = getOperateData(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = service.CheckLogin(r); err != nil {
		http.Error(w, fmt.Sprintf("未授权: %v", err), http.StatusUnauthorized)
		return
	}
	if resData, err = operate.DeleteOperate(dataOperate); err != nil {
		return
	}
	w.Header().Set("Content-Type", resData.ContentType)
	w.WriteHeader(resData.Header)
	if err = json.NewEncoder(w).Encode(resData.ResData); err != nil {
		http.Error(w, "服务器错误", http.StatusInternalServerError)
	}

}

func HandleQuery(w http.ResponseWriter, r *http.Request) {
	var (
		err         error
		resData     *operate.QueryResponse
		dataOperate map[string][]string
	)
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		http.Error(w, "Unsupported Content-Type, must be application/x-www-form-urlencoded", http.StatusUnsupportedMediaType)
		return
	}
	if dataOperate, err = getOperateData(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = service.CheckLogin(r); err != nil {
		http.Error(w, fmt.Sprintf("未授权: %v", err), http.StatusUnauthorized)
		return
	}
	if resData, err = operate.QueryOperate(dataOperate); err != nil {
		return
	}
	w.Header().Set("Content-Type", resData.ContentType)
	w.WriteHeader(resData.Header)
	if err = json.NewEncoder(w).Encode(resData.ResData); err != nil {
		http.Error(w, "服务器错误", http.StatusInternalServerError)
	}
}

func HandleAmount(w http.ResponseWriter, r *http.Request) {
	var (
		err         error
		resData     *operate.AmountResponse
		dataOperate map[string][]string
	)
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if dataOperate, err = getOperateData(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = service.CheckLogin(r); err != nil {
		http.Error(w, fmt.Sprintf("未授权: %v", err), http.StatusUnauthorized)
		return
	}
	if resData, err = operate.AmountOperate(dataOperate); err != nil {
		return
	}
	w.Header().Set("Content-Type", resData.ContentType)
	w.WriteHeader(resData.Header)
	if err = json.NewEncoder(w).Encode(resData.ResData); err != nil {
		http.Error(w, "服务器错误", http.StatusInternalServerError)
	}
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var (
		err         error
		resData     *operate.LoginResponse
		dataOperate map[string][]string
	)
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		http.Error(w, "Unsupported Content-Type, must be application/x-www-form-urlencoded", http.StatusUnsupportedMediaType)
		return
	}
	if dataOperate, err = getOperateData(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if resData, err = operate.LoginOperate(dataOperate); err != nil {
		return
	}
	w.Header().Set("Content-Type", resData.ContentType)
	w.WriteHeader(resData.Header)
	if err = json.NewEncoder(w).Encode(resData.ResData); err != nil {
		http.Error(w, "服务器错误", http.StatusInternalServerError)
	}
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	var (
		err         error
		resData     *operate.RegisterResponse
		dataOperate map[string][]string
	)
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		http.Error(w, "Unsupported Content-Type, must be application/x-www-form-urlencoded", http.StatusUnsupportedMediaType)
		return
	}
	if dataOperate, err = getOperateData(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if resData, err = operate.HandleRegister(dataOperate); err != nil {
		return
	}
	w.Header().Set("Content-Type", resData.ContentType)
	w.WriteHeader(resData.Header)
	if err = json.NewEncoder(w).Encode(resData.ResData); err != nil {
		http.Error(w, "服务器错误", http.StatusInternalServerError)
	}
}

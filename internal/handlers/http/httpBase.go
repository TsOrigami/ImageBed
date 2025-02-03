package http

import (
	"ImageV2/internal/handlers/operate"
	service "ImageV2/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func HandleHttp() {
	http.HandleFunc("/upload", HandleUpload)                // 上传图片
	http.HandleFunc("/invoke/", operate.HandleInvoke)       // 调用图片
	http.HandleFunc("/delete", HandleDelete)                // 删除图片
	http.HandleFunc("/query", HandleQuery)                  // 查询原图
	http.HandleFunc("/thumbnail/", operate.HandleThumbnail) // 查询缩略图
	http.HandleFunc("/login", HandleLogin)                  // 登录
}

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	var (
		err error
	)
	dataOperate := make(map[string][]string)
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Printf("Content-Type: %s\n", r.Header.Get("Content-Type"))
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

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var (
		err         error
		resData     *operate.LoginResponse
		dataOperate map[string][]string
	)
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

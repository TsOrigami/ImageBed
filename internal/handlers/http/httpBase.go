package http

import (
	"ImageV2/internal/handlers/operate"
	service "ImageV2/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
)

func HandleHttp() {
	http.HandleFunc("/upload", operate.HandleUpload)  // 上传图片
	http.HandleFunc("/invoke/", operate.HandleInvoke) // 调用图片
	http.HandleFunc("/delete", HandleDelete)          // 删除图片
	http.HandleFunc("/query", HandleQuery)            // 查询图片
	http.HandleFunc("/login", HandleLogin)            // 登录
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
		http.Error(w, "Unsupported Content-Type", http.StatusUnsupportedMediaType)
		return
	}
	// 检查登录状态
	err = service.CheckLogin(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("未授权: %v", err), http.StatusUnauthorized)
		return
	}
	// 解析表单数据
	if err = r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("解析表单数据失败: %v", err), http.StatusBadRequest)
		return
	}
	for key, values := range r.Form {
		dataOperate[key] = append(dataOperate[key], values...)
	}
	for key, values := range r.Header {
		dataOperate[key] = append(dataOperate[key], values...)
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
	// 确保 Content-Type 是 application/x-www-form-urlencoded
	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		http.Error(w, "Unsupported Content-Type", http.StatusUnsupportedMediaType)
		return
	}
	// 检查登录状态
	err = service.CheckLogin(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("未授权: %v", err), http.StatusUnauthorized)
		return
	}
	// 解析表单数据
	if err = r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("解析表单数据失败: %v", err), http.StatusBadRequest)
		return
	}
	for key, values := range r.Form {
		dataOperate[key] = append(dataOperate[key], values...)
	}
	for key, values := range r.Header {
		dataOperate[key] = append(dataOperate[key], values...)
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
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	// 确保 Content-Type 是 application/x-www-form-urlencoded
	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		http.Error(w, "Unsupported Content-Type", http.StatusUnsupportedMediaType)
		return
	}
	if err = r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("解析表单数据失败: %v", err), http.StatusBadRequest)
		return
	}
	for key, values := range r.Form {
		dataOperate[key] = append(dataOperate[key], values...)
	}
	for key, values := range r.Header {
		dataOperate[key] = append(dataOperate[key], values...)
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

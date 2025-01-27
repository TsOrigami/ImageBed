package handlers

import (
	dbImage "ImageV2/internal/db/image"
	dbUser "ImageV2/internal/db/user"
	service "ImageV2/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func HandleQuery(w http.ResponseWriter, r *http.Request) {
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
	err := service.CheckLogin(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("未授权: %v", err), http.StatusUnauthorized)
		return
	}
	// 解析表单数据
	err = r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("解析表单数据失败: %v", err), http.StatusBadRequest)
		return
	}
	token := r.Header.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	}
	username, err := dbUser.GetUsername(token)
	startStr := r.Form.Get("start")
	endStr := r.Form.Get("end")
	uuids, err := dbImage.GetInfoQuery(startStr, endStr, username)
	if err != nil {
		return
	}
	response := ImageQueryResponse{
		Code:  200,   // 设置 code 为 200
		UUIDs: uuids, // 赋值查询到的 uuids
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "服务器错误", http.StatusInternalServerError)
	}
}

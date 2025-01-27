package handlers

import (
	service "ImageV2/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	// 确保 Content-Type 是 application/x-www-form-urlencoded
	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		http.Error(w, "Unsupported Content-Type", http.StatusUnsupportedMediaType)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("解析表单数据失败: %v", err), http.StatusBadRequest)
		return
	}
	// 获取表单中的 uuid 参数
	account := r.Form.Get("account")
	salt := r.Form.Get("salt")
	sign := r.Form.Get("sign")
	token, userName, err := service.Login(account, salt, sign)
	if err != nil {
		http.Error(w, fmt.Sprintf("登录失败: %v", err), http.StatusBadRequest)
		return
	}
	response := UserResponse{
		Token:    token,
		Username: userName,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "服务器错误", http.StatusInternalServerError)
	}
}

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

func HandleDelete(w http.ResponseWriter, r *http.Request) {
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
	// 获取表单中的 uuid 参数
	token := r.Header.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	}
	usernameDel, err := dbUser.GetUsername(token)
	uuid := r.Form.Get("uuid")
	err = dbImage.DeleteInfoFromSQL(uuid, usernameDel)
	if err != nil {
		_, err := fmt.Fprintf(w, "删除数据失败: %v", err)
		if err != nil {
			return
		}
		return
	}
	response := ImageResponse{
		Code: 200,
		Msg:  "删除成功:" + uuid,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "服务器错误", http.StatusInternalServerError)
	}
}

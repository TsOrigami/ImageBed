package handlers

import (
	dbImage "ImageV2/internal/db/image"
	service "ImageV2/internal/services"
	"fmt"
	"net/http"
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
	uuid := r.Form.Get("uuid")
	err = dbImage.DeleteInfoFromSQL(uuid)
	if err != nil {
		_, err := fmt.Fprintf(w, "删除数据失败: %v", err)
		if err != nil {
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("删除成功:" + uuid))
	if err != nil {
		return
	}
}

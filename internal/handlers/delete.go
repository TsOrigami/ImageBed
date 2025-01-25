package handlers

import (
	database "ImageV2/internal/db"
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
	// 解析表单数据
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("解析表单数据失败: %v", err), http.StatusBadRequest)
		return
	}
	// 获取表单中的 uuid 参数
	uuid := r.Form.Get("uuid")
	err = database.DeleteInfoFromSQL(uuid)
	if err != nil {
		http.Error(w, "删除数据失败", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("删除成功:" + uuid))
	if err != nil {
		return
	}
}

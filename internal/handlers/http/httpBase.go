package http

import (
	dbUser "ImageV2/internal/db/user"
	"ImageV2/internal/handlers/operate"
	service "ImageV2/internal/services"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"
)

func HandleHttp() {
	http.HandleFunc("/upload", HandleUpload)          // 上传图片
	http.HandleFunc("/invoke/", operate.HandleInvoke) // 调用图片
	http.HandleFunc("/delete", HandleDelete)          // 删除图片
	http.HandleFunc("/query", HandleQuery)            // 查询图片
	http.HandleFunc("/login", HandleLogin)            // 登录
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

func getOperateFile(r *http.Request) (map[string][]string, error) {
	var (
		err      error
		token    string
		username string
	)
	dataOperate := make(map[string][]string)
	token = r.Header.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	}
	if username, err = dbUser.GetUsername(token); err != nil {
		return nil, err
	}
	dataOperate["username"] = append(dataOperate["username"], username)
	if err = r.ParseMultipartForm(100 << 20); err != nil {
		return nil, err
	}
	mForm := r.MultipartForm
	var wg sync.WaitGroup
	for k := range mForm.File {
		wg.Add(1) // 增加计数
		go func(k string) {
			defer wg.Done() // 完成时减少计数
			file, fileHeader, err := r.FormFile(k)
			if err != nil {
				return
			}
			defer func(file multipart.File) {
				err := file.Close()
				if err != nil {
					return
				}
			}(file)
			fileName := fileHeader.Filename
			fileBytes := make([]byte, fileHeader.Size)
			if _, err = file.Read(fileBytes); err != nil {
				return
			}
			encodedFile := base64.StdEncoding.EncodeToString(fileBytes)
			dataOperate["pic_"+fileName] = append(dataOperate["pic_"+fileName], encodedFile)
		}(k)
	}
	wg.Wait() // 等待所有文件上传任务完成
	return dataOperate, nil
}

func getOperateData(r *http.Request) (map[string][]string, error) {
	var (
		err      error
		token    string
		username string
	)
	dataOperate := make(map[string][]string)
	if r.Method != http.MethodPost {
		return nil, fmt.Errorf("method Not Allowed")
	}
	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		return nil, fmt.Errorf("unsupported Content-Type")
	}
	if err = r.ParseForm(); err != nil {
		return nil, fmt.Errorf("解析表单数据失败: %v", err)
	}
	token = r.Header.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	}
	if username, err = dbUser.GetUsername(token); err == nil {
		dataOperate["username"] = append(dataOperate["username"], username)
	}
	for key, values := range r.Form {
		dataOperate[key] = append(dataOperate[key], values...)
	}
	for key, values := range r.Header {
		dataOperate[key] = append(dataOperate[key], values...)
	}
	return dataOperate, nil
}

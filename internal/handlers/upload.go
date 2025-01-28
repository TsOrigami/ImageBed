package handlers

import (
	dbImage "ImageV2/internal/db/image"
	dbUser "ImageV2/internal/db/user"
	errorHandle "ImageV2/internal/error"
	service "ImageV2/internal/services"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"
	"time"
)

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		username string
		token    string
	)
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Printf("Content-Type: %s\n", r.Header.Get("Content-Type"))
	// 检查 Content-Type 是否以 "multipart/form-data" 开头
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		http.Error(w, "Unsupported Content-Type, must be multipart/form-data", http.StatusUnsupportedMediaType)
		return
	}
	// 检查登录状态
	if err = service.CheckLogin(r); err != nil {
		http.Error(w, fmt.Sprintf("未授权: %v", err), http.StatusUnauthorized)
		return
	}
	token = r.Header.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	}
	if username, err = dbUser.GetUsername(token); err != nil {
		errorHandle.DatabaseError(w, err)
		return
	}
	if err = r.ParseMultipartForm(100 << 20); err != nil {
		errorHandle.UploadError(w)
		return
	} // 限制上传文件大小为100MB
	mForm := r.MultipartForm
	var wg sync.WaitGroup

	// 逐个文件并发处理
	for k := range mForm.File {
		wg.Add(1) // 增加计数
		go func(k string) {
			defer wg.Done() // 完成时减少计数
			file, fileHeader, err := r.FormFile(k)
			if err != nil {
				errorHandle.UploadError(w)
				return
			}
			defer func(file multipart.File) {
				err := file.Close()
				if err != nil {
					return
				}
			}(file)
			fileName := fileHeader.Filename
			imagePath, err := service.GetSavePath()
			if err != nil {
				errorHandle.UploadError(w)
				return
			}
			err = service.SaveImage(imagePath, fileName, file)
			if err != nil {
				errorHandle.UploadError(w)
				return
			}
			localFileName := imagePath + "/" + fileName
			if err != nil {
				errorHandle.UploadError(w)
				return
			}
			channel := make(chan string, 1)
			go func() {
				picSha256, err := service.GetSha256(localFileName)
				if err != nil {
					errorHandle.UploadError(w)
					return
				}
				channel <- picSha256
			}()
			picSha256 := <-channel
			uploadTime := time.Now()
			err = dbImage.SaveInfoToSQL(fileName, username, picSha256, uploadTime)
			if err != nil {
				errorHandle.DatabaseError(w, err)
				return
			}
		}(k)
	}
	wg.Wait() // 等待所有文件上传任务完成
	response := ImageResponse{
		Code: 200,
		Msg:  "文件上传成功",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "服务器错误", http.StatusInternalServerError)
	}
}

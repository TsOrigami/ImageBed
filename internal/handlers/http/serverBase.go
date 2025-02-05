package http

import (
	dbUser "ImageV2/internal/db/user"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"
)

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
		if dataOperate[key] == nil {
			dataOperate[key] = append(dataOperate[key], values...)
		}
	}
	for key, values := range r.Header {
		if dataOperate[key] == nil {
			dataOperate[key] = append(dataOperate[key], values...)
		}
	}
	return dataOperate, nil
}

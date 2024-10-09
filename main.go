package main

import (
	sql "ImageBed_api/SQL_Operate"
	api "ImageBed_api/api"
	errorHandle "ImageBed_api/error"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const imagesPath = "./images"

// 计算图片的MD5值
func getMd5(FilePath string) string {
	md5file, err := os.Open(FilePath)
	if err != nil {
		return ""
	}
	hash := md5.New()
	_, _ = io.Copy(hash, md5file)
	picMd5 := hex.EncodeToString(hash.Sum(nil))
	return picMd5
}

func handleUploadImages(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(100 << 20) // 限制上传文件大小为100MB
	mForm := r.MultipartForm

	for k := range mForm.File {
		file, fileHeader, err := r.FormFile(k)
		if err != nil {
			errorHandle.UploadError(w)
			return
		}
		defer file.Close()

		fileName := fileHeader.Filename
		localFileName := imagesPath + "/" + fileName
		out, err := os.Create(localFileName)
		if err != nil {
			errorHandle.UploadError(w)
			return
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			errorHandle.UploadError(w)
			return
		}

		picMd5 := getMd5(localFileName)

		uploadTime := time.Now()

		err = sql.SaveImageInfo(fileName, picMd5, uploadTime)
		if err != nil {
			errorHandle.DatabaseError(w, err)
			return
		}

		fmt.Println("md5:", picMd5)
		fmt.Printf("file %s uploaded ok\n", fileHeader.Filename)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("文件上传成功"))
	}
}

func handlerDeleteImages(w http.ResponseWriter, r *http.Request) {

	type DeleteRequest struct {
		MD5 string `json:"md5"`
	}

	if r.Method != http.MethodDelete {
		errorHandle.MethodNotAllowedError(w)
		return
	}

	// 解析JSON请求体
	var req DeleteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorHandle.BadRequestError(w)
		return
	}

	// 调用删除函数
	err = sql.DeleteImageByMD5(req.MD5)
	if err != nil {
		errorHandle.InvalidMD5Error(w)
		return
	}

	// 返回成功响应
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("删除成功"))
}

func handlerInquireImages(w http.ResponseWriter, r *http.Request) {
	// 确保是GET请求
	if r.Method != http.MethodGet {
		errorHandle.MethodNotAllowedError(w)
		return
	}

	// 获取start和end参数
	startParam := r.URL.Query().Get("start")
	endParam := r.URL.Query().Get("end")

	// 将参数转换为整数
	start, err := strconv.Atoi(startParam)
	if err != nil {
		errorHandle.InvalidParameterError(w, "start")
		return
	}

	end, err := strconv.Atoi(endParam)
	if err != nil {
		errorHandle.InvalidParameterError(w, "end")
		return
	}

	// 调用函数获取数据
	images, err := sql.GetImageInfo(start, end)
	if err != nil {
		errorHandle.InternalServerError(w, err)
		return
	}

	// 将结果编码为JSON并返回
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(images)

}

func handlerFetchImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorHandle.MethodNotAllowedError(w)
		return
	}

	// 从查询参数中获取MD5值
	md5Hash := r.URL.Query().Get("md5")
	if md5Hash == "" {
		errorHandle.InvalidMD5Error(w)
		return
	}

	// 根据MD5值获取图片
	filePath, err := sql.FetchImageByMD5(md5Hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("无法打开文件: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, fmt.Sprintf("获取文件信息失败: %v", err), http.StatusInternalServerError)
		return
	}

	// 根据文件扩展名设置响应头
	ext := filepath.Ext(filePath) // 获取文件扩展名
	var contentType string
	switch ext {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	default:
		http.Error(w, "不支持的文件类型", http.StatusUnsupportedMediaType)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Last-Modified", fileInfo.ModTime().Format(http.TimeFormat))
	w.WriteHeader(http.StatusOK)

	// 将文件内容返回给客户端
	http.ServeContent(w, r, filePath, fileInfo.ModTime(), file)
}

func main() {
	http.HandleFunc("/upload", handleUploadImages)       // 上传图片
	http.HandleFunc("/delete", handlerDeleteImages)      // 删除图片
	http.HandleFunc("/inquire", handlerInquireImages)    // 查询位于范围内的图片信息 (名字、路径、MD5值、创建时间)
	http.HandleFunc("/fetch", handlerFetchImage)         // get请求获取图片
	http.HandleFunc("/image/api/", api.HandleImageByMd5) // 获取图片
	err := http.ListenAndServe(":8000", nil)             // 启动服务器
	if err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
		return
	}
}

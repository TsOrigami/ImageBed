package handlers

import (
	database "ImageV2/internal/db"
	errorHandle "ImageV2/internal/error"
	"ImageV2/internal/services"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		errorHandle.UploadError(w)
		return
	} // 限制上传文件大小为100MB

	mForm := r.MultipartForm
	for k := range mForm.File {
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
		imagePath, err := services.GetSavePath()
		if err != nil {
			errorHandle.UploadError(w)
			return
		}
		localFileName := imagePath + "/" + fileName
		out, err := os.Create(localFileName)
		if err != nil {
			errorHandle.UploadError(w)
			return
		}
		defer func(out *os.File) {
			err := out.Close()
			if err != nil {
				return
			}
		}(out)

		_, err = io.Copy(out, file)
		if err != nil {
			errorHandle.UploadError(w)
			return
		}

		picSha256, err := services.GetSha256(localFileName)
		if err != nil {
			errorHandle.UploadError(w)
			return
		}
		uploadTime := time.Now()

		err = SaveImageInfo(fileName, picSha256, uploadTime)
		if err != nil {
			errorHandle.DatabaseError(w, err)
			return
		}

		fmt.Println("sha256:", picSha256)
		fmt.Printf("file %s uploaded ok\n", fileHeader.Filename)
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("文件上传成功"))
	if err != nil {
		return
	}
}

func SaveImageInfo(imageName, sha256Hash string, createdAt time.Time) error {
	// 获取数据库连接
	dbInfo, err := database.GetDB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %v", err)
	}

	// 插入数据
	insertSQL := `INSERT INTO image_info (image_name, sha256Hash, created_at) VALUES (?, ?, ?)`
	_, err = dbInfo.Connect.Exec(insertSQL, imageName, sha256Hash, createdAt)
	if err != nil {
		fmt.Println(sha256Hash)
		return fmt.Errorf("插入数据失败: %v", err)
	}

	fmt.Printf("成功插入图片信息: image_name=%s, sha256_hash=%s, created_at=%s\n", imageName, sha256Hash, createdAt)
	return nil
}

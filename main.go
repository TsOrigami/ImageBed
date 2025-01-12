package main

import (
	errorHandle "ImageV2/error"
	dbcontrol "ImageV2/sqlcontrol"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

const imagesPath = "./images"

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

func SaveImageInfo(imageName, md5Hash string, createdAt time.Time) error {
	// 获取数据库连接
	db, err := dbcontrol.ConnectionDB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %v", err)
	}

	// 创建表（如果不存在）
	createTableSQL := `CREATE TABLE IF NOT EXISTS image_info (
		image_name VARCHAR(255),
		md5_hash CHAR(32) PRIMARY KEY,
		created_at DATETIME
	)`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("创建表失败: %v", err)
	}

	// 插入数据
	insertSQL := `INSERT INTO image_info (image_name, md5_hash, created_at) VALUES (?, ?, ?)`
	_, err = db.Exec(insertSQL, imageName, md5Hash, createdAt)
	if err != nil {
		return fmt.Errorf("插入数据失败: %v", err)
	}

	fmt.Printf("成功插入图片信息: image_name=%s, md5_hash=%s, created_at=%s\n", imageName, md5Hash, createdAt)
	return nil
}

func handleUploadImages(w http.ResponseWriter, r *http.Request) {
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
		localFileName := imagesPath + "/" + fileName
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

		picMd5 := getMd5(localFileName)
		uploadTime := time.Now()

		err = SaveImageInfo(fileName, picMd5, uploadTime)
		if err != nil {
			errorHandle.DatabaseError(w, err)
			return
		}

		fmt.Println("md5:", picMd5)
		fmt.Printf("file %s uploaded ok\n", fileHeader.Filename)
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("文件上传成功"))
	if err != nil {
		return
	}
}

func main() {
	http.HandleFunc("/upload", handleUploadImages) // 上传图片
	err := http.ListenAndServe(":8000", nil)       // 启动服务器
	if err != nil {
		fmt.Printf("服务器启动失败: %v\n", err)
		return
	}
}

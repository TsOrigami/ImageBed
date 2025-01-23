package handlers

import (
	database "ImageV2/internal/db"
	errorHandle "ImageV2/internal/error"
	"ImageV2/internal/services"
	"net/http"
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
		err = file.Close()
		if err != nil {
			return
		}
		fileName := fileHeader.Filename
		imagePath, err := services.GetSavePath()
		err = services.SaveImage(imagePath, fileName, file)
		if err != nil {
			return
		}
		localFileName := imagePath + "/" + fileName
		picSha256, err := services.GetSha256(localFileName)
		if err != nil {
			errorHandle.UploadError(w)
			return
		}
		uploadTime := time.Now()
		err = database.SaveInfoToSQL(fileName, picSha256, uploadTime)
		if err != nil {
			errorHandle.DatabaseError(w, err)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("文件上传成功"))
	if err != nil {
		return
	}
}

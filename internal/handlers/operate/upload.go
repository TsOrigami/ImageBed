package operate

import (
	dbImage "ImageV2/internal/db/image"
	"ImageV2/internal/handlers"
	service "ImageV2/internal/services"
	"encoding/base64"
	"net/http"
	"sync"
	"time"
)

type UploadResponse struct {
	ResData     *handlers.ImageResponse `json:"ResData"`
	ContentType string                  `json:"Content-Type"`
	Header      int                     `json:"Header"`
}

func UploadOperate(dataOperate map[string][]string) (*UploadResponse, error) {
	var (
		err       error
		username  string
		imagePath string
	)
	if imagePath, err = service.GetSavePath(); err != nil {
		return nil, err
	}
	errChan := make(chan error, len(dataOperate))
	username = dataOperate["username"][0]
	var wg sync.WaitGroup
	for key, values := range dataOperate {
		wg.Add(1)
		if key[:3] != "pic" {
			wg.Done()
			continue
		}
		go func(key string, values []string) {
			defer wg.Done()
			filename := key[4:]
			fileBytes, err := base64.StdEncoding.DecodeString(values[0])
			localFileName := imagePath + "/" + filename
			if err = service.SaveImage(imagePath, filename, fileBytes); err != nil {
				errChan <- err
				return
			}
			channel := make(chan string, 1)
			go func() {
				picSha256, err := service.GetSha256(localFileName)
				if err != nil {
					errChan <- err
					return
				}
				channel <- picSha256
			}()
			picSha256 := <-channel
			uploadTime := time.Now()
			if err = dbImage.SaveInfoToSQL(filename, username, picSha256, uploadTime); err != nil {
				return
			}
		}(key, values)
	}
	go func() {
		wg.Wait()
		close(errChan)
	}()
	if err = <-errChan; err != nil {
		return nil, err
	}
	var response = UploadResponse{
		ResData: &handlers.ImageResponse{
			Code: 200,
			Msg:  "文件上传成功",
		},
		ContentType: "application/json",
		Header:      http.StatusOK,
	}
	return &response, nil
}

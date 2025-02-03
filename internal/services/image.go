package services

import (
	conf "ImageV2/configs"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/disintegration/imaging"
	"io"
	"os"
	"sync"
)

// SaveImage 保存图片
func SaveImage(imagePath string, ThumbnailPath string, fileName string, file []byte) error {
	localFileName := imagePath + "/" + fileName
	thumbnailFileName := ThumbnailPath + "/" + fileName
	out, err := os.Create(localFileName)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			return
		}
	}(out)
	reader := bytes.NewReader(file)
	_, err = io.Copy(out, reader)
	if err != nil {
		return err
	}
	//生成缩略图
	img, err := imaging.Decode(bytes.NewReader(file))
	if err != nil {
		return err
	}
	thumbnail := imaging.Resize(img, 0, 240, imaging.Lanczos)
	err = imaging.Save(thumbnail, thumbnailFileName)
	if err != nil {
		return err
	}
	return nil
}

var (
	ImagePath     = ""
	ThumbnailPath = ""
	mu            sync.Mutex
)

func GetSavePath() (string, string, error) {
	if ImagePath == "" {
		mu.Lock()
		defer mu.Unlock()
		if ImagePath == "" {
			var err error
			var jsonData []byte
			jsonData, err = conf.GetConfigGroupAsJSON("server")
			if err != nil {
				return "", "", fmt.Errorf("获取配置信息失败: %v", err)
			}
			var config map[string]string
			err = json.Unmarshal(jsonData, &config)
			if err != nil {
				return "", "", fmt.Errorf("解析配置信息失败: %v", err)
			}
			ImagePath = config["imagesPath"]
			ThumbnailPath = config["thumbnailPath"]
		}
	}
	return ImagePath, ThumbnailPath, nil
}

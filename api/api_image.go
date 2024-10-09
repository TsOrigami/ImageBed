package api

import (
	sql "ImageBed_api/SQL_Operate"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func HandleImageByMd5(w http.ResponseWriter, r *http.Request) {
	// 从URL中提取MD5值
	md5Value := strings.TrimPrefix(r.URL.Path, "/image/api/") // 去掉前缀以获取MD5值

	if md5Value == "" {
		http.Error(w, "MD5 值不能为空", http.StatusBadRequest)
		return
	}

	// 根据MD5值查询数据库获取对应的image_name
	imageInfo, err := sql.GetImageInfoByMD5(md5Value)
	if err != nil {
		http.Error(w, fmt.Sprintf("查询数据库失败: %v", err), http.StatusInternalServerError)
		return
	}

	// 拼接图片的完整路径
	imageFilePath := sql.ImagesPath + "/" + imageInfo.ImageName

	// 检查图片文件是否存在
	if _, err := os.Stat(imageFilePath); os.IsNotExist(err) {
		http.Error(w, "图片不存在", http.StatusNotFound)
		return
	}

	// 打开图片文件
	file, err := os.Open(imageFilePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("无法打开图片文件: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 根据文件扩展名设置 MIME 类型
	ext := strings.ToLower(filepath.Ext(imageFilePath))
	switch ext {
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	default:
		http.Error(w, "不支持的图片格式", http.StatusUnsupportedMediaType)
		return
	}

	// 将图片内容写入 HTTP 响应
	http.ServeFile(w, r, imageFilePath)
}

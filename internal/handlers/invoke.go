package handlers

import (
	database "ImageV2/internal/db"
	"ImageV2/internal/services"
	"fmt"
	"net/http"
	"regexp"
)

func HandleInvoke(w http.ResponseWriter, r *http.Request) {
	re := regexp.MustCompile(`/invoke/([a-fA-F0-9-]+)`)
	matches := re.FindStringSubmatch(r.URL.Path)
	if len(matches) > 1 {
		picInfo, err := database.GetInfoByUUID(matches[1])
		if err != nil {
			http.Error(w, fmt.Sprintf("获取数据失败: %v", err), http.StatusInternalServerError)
			return
		}
		var ImagePath string
		ImagePath, err = services.GetSavePath()
		if err != nil {
			http.Error(w, fmt.Sprintf("获取图片路径失败: %v", err), http.StatusInternalServerError)
			return
		}
		ImagePath = ImagePath + "/" + picInfo.ImageName
		w.Header().Set("Content-Type", "image/jpeg")
		http.ServeFile(w, r, ImagePath)
		return
	} else {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
}

package operate

import (
	dbImage "ImageV2/internal/db/image"
	"ImageV2/internal/services"
	"fmt"
	"net/http"
	"regexp"
)

func HandleThumbnail(w http.ResponseWriter, r *http.Request) {
	re := regexp.MustCompile(`/thumbnail/([a-fA-F0-9-]+)`)
	matches := re.FindStringSubmatch(r.URL.Path)
	if len(matches) > 1 {
		picInfo, err := dbImage.GetInfoByUUID(matches[1])
		if err != nil {
			http.Error(w, fmt.Sprintf("获取数据失败: %v", err), http.StatusInternalServerError)
			return
		}
		var ThumbnailPath string
		_, ThumbnailPath, err = services.GetSavePath()
		if err != nil {
			http.Error(w, fmt.Sprintf("获取图片路径失败: %v", err), http.StatusInternalServerError)
			return
		}
		ThumbnailPath = ThumbnailPath + "/" + picInfo.ImageName
		w.Header().Set("Content-Type", "image/jpeg")
		http.ServeFile(w, r, ThumbnailPath)
		return
	} else {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
}

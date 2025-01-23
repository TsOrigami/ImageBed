package handlers

import (
	database "ImageV2/internal/db"
	"net/http"
	"regexp"
)

func HandleDelete(w http.ResponseWriter, r *http.Request) {
	re := regexp.MustCompile(`/delete/([a-fA-F0-9-]+)`)
	matches := re.FindStringSubmatch(r.URL.Path)
	if len(matches) > 1 {
		err := database.DeleteInfoFromSQL(matches[1])
		if err != nil {
			http.Error(w, "删除数据失败", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
}

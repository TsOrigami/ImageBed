package operate

import (
	conf "ImageV2/configs"
	"ImageV2/internal/handlers"
	"encoding/json"
	"net/http"
)

type registerResponse struct {
	ResData     *handlers.SystemRegister `json:"ResData"`
	ContentType string                   `json:"Content-Type"`
	Header      int                      `json:"Header"`
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		response registerResponse
		jsonData []byte
	)
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	jsonData, err = conf.GetConfigGroupAsJSON("server")
	if err != nil {
		return
	}
	var systemConfig = make(map[string]string)
	err = json.Unmarshal(jsonData, &systemConfig)
	prefixPath := systemConfig["prefix"]
	response.ResData = &handlers.SystemRegister{
		Register: prefixPath,
	}
	response.ContentType = "application/json"
	response.Header = http.StatusOK
	w.Header().Set("Content-Type", response.ContentType)
	w.WriteHeader(response.Header)
	err = json.NewEncoder(w).Encode(response.ResData)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	return
}

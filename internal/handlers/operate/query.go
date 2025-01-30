package operate

import (
	dbImage "ImageV2/internal/db/image"
	"ImageV2/internal/handlers"
	"net/http"
)

type QueryResponse struct {
	ResData     *handlers.ImageQueryResponse `json:"ResData"`
	ContentType string                       `json:"Content-Type"`
	Header      int                          `json:"Header"`
}

func QueryOperate(dataOperate map[string][]string) (*QueryResponse, error) {
	var (
		err      error
		username string
		uuids    []string
		startStr string
		endStr   string
		response QueryResponse
	)
	username = dataOperate["username"][0]
	startStr = dataOperate["start"][0]
	endStr = dataOperate["end"][0]
	uuids, err = dbImage.GetInfoQuery(startStr, endStr, username)
	if err != nil {
		return nil, err
	}
	response = QueryResponse{
		ResData: &handlers.ImageQueryResponse{
			Code:  200,   // 设置 code 为 200
			UUIDs: uuids, // 赋值查询到的 uuids
		},
		ContentType: "application/json",
		Header:      http.StatusOK,
	}
	return &response, nil
}

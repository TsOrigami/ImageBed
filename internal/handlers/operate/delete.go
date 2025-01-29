package operate

import (
	dbImage "ImageV2/internal/db/image"
	dbUser "ImageV2/internal/db/user"
	"ImageV2/internal/handlers"
	"net/http"
	"strings"
)

type DeleteResponse struct {
	ResData     *handlers.ImageResponse `json:"ResData"`
	ContentType string                  `json:"Content-Type"`
	Header      int                     `json:"Header"`
}

func DeleteOperate(dataOperate map[string][]string) (*DeleteResponse, error) {
	// 获取表单中的 uuid 参数
	token := dataOperate["Authorization"][0]
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	}
	usernameDel, err := dbUser.GetUsername(token)
	uuid := dataOperate["uuid"][0]
	err = dbImage.DeleteInfoFromSQL(uuid, usernameDel)
	if err != nil {
		return nil, err
	}
	var response = DeleteResponse{
		ResData: &handlers.ImageResponse{
			Code: 200,
			Msg:  "删除成功:" + uuid,
		},
		ContentType: "application/json",
		Header:      http.StatusOK,
	}
	return &response, nil
}

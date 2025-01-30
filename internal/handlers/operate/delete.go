package operate

import (
	dbImage "ImageV2/internal/db/image"
	"ImageV2/internal/handlers"
	"net/http"
)

type DeleteResponse struct {
	ResData     *handlers.ImageResponse `json:"ResData"`
	ContentType string                  `json:"Content-Type"`
	Header      int                     `json:"Header"`
}

func DeleteOperate(dataOperate map[string][]string) (*DeleteResponse, error) {
	// 获取表单中的 uuid 参数
	var (
		err         error
		usernameDel string
		uuid        string
	)
	usernameDel = dataOperate["username"][0]
	uuid = dataOperate["uuid"][0]
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

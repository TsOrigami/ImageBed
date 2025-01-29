package operate

import (
	"ImageV2/internal/handlers"
	service "ImageV2/internal/services"
	"net/http"
)

type LoginResponse struct {
	ResData     *handlers.UserResponse `json:"ResData"`
	ContentType string                 `json:"Content-Type"`
	Header      int                    `json:"Header"`
}

func LoginOperate(dataOperate map[string][]string) (*LoginResponse, error) {
	var (
		err      error
		account  string
		salt     string
		sign     string
		token    string
		userName string
		response LoginResponse
	)
	// 获取表单中的 uuid 参数
	account = dataOperate["account"][0]
	salt = dataOperate["salt"][0]
	sign = dataOperate["sign"][0]
	if token, userName, err = service.Login(account, salt, sign); err != nil {
		return nil, err
	}
	response = LoginResponse{
		ResData: &handlers.UserResponse{
			Token:    token,
			Username: userName,
		},
		ContentType: "application/json",
		Header:      http.StatusOK,
	}
	return &response, nil
}

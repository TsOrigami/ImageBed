package operate

import (
	"ImageV2/internal/handlers"
	service "ImageV2/internal/services"
	"fmt"
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
	fmt.Printf("dataOperate: %v\n", dataOperate)
	if dataOperate["account"] != nil && dataOperate["account"][0] != "" {
		account = dataOperate["account"][0]
	} else {
		err = fmt.Errorf("account is empty")
	}
	if dataOperate["salt"] != nil && dataOperate["salt"][0] != "" {
		salt = dataOperate["salt"][0]
	} else {
		err = fmt.Errorf("salt is empty")
	}
	if dataOperate["sign"] != nil || dataOperate["sign"][0] != "" {
		sign = dataOperate["sign"][0]
	} else {
		err = fmt.Errorf("sign is empty")
	}
	if err != nil {
		return nil, err
	}
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

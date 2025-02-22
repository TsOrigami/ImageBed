package operate

import (
	"ImageV2/internal/handlers"
	service "ImageV2/internal/services"
	"fmt"
	"net/http"
)

type RegisterResponse struct {
	ResData     *handlers.SystemRegister `json:"ResData"`
	ContentType string                   `json:"Content-Type"`
	Header      int                      `json:"Header"`
}

func HandleRegister(dataOperate map[string][]string) (*RegisterResponse, error) {
	var (
		err      error
		account  string
		username string
		password string
	)
	if dataOperate["account"] != nil && dataOperate["account"][0] != "" {
		account = dataOperate["account"][0]
	} else {
		err = fmt.Errorf("account is empty")
	}
	if dataOperate["user"] != nil && dataOperate["user"][0] != "" {
		username = dataOperate["user"][0]
	} else {
		err = fmt.Errorf("username is empty")
	}
	if dataOperate["password"] != nil && dataOperate["password"][0] != "" {
		password = dataOperate["password"][0]
	} else {
		err = fmt.Errorf("password is empty")
	}
	if err != nil {
		return nil, err
	}
	err = service.Registered(account, username, password)
	if err != nil {
		return nil, err
	}
	response := RegisterResponse{
		ResData: &handlers.SystemRegister{
			Register: "success",
		},
		ContentType: "application/json",
		Header:      http.StatusOK,
	}
	return &response, nil
}

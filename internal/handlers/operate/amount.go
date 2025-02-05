package operate

import (
	dbUser "ImageV2/internal/db/user"
	"ImageV2/internal/handlers"
	"net/http"
)

type AmountResponse struct {
	ResData     *handlers.ImageAmountResponse `json:"ResData"`
	ContentType string                        `json:"Content-Type"`
	Header      int                           `json:"Header"`
}

func AmountOperate(dataOperate map[string][]string) (*AmountResponse, error) {
	var (
		err      error
		username string
		count    int
		response AmountResponse
	)
	username = dataOperate["username"][0]
	count, err = dbUser.GetAmount(username)
	if err != nil {
		return nil, err
	}
	response = AmountResponse{
		ResData: &handlers.ImageAmountResponse{
			Code:  200,   // 设置 code 为 200
			Count: count, // 赋值查询到的 uuids
		},
		ContentType: "application/json",
		Header:      http.StatusOK,
	}
	return &response, nil
}

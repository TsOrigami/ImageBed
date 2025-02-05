package handlers

type ImageResponse struct {
	Code  int      `json:"code"`
	Msg   string   `json:"msg"`
	UUIDs []string `json:"uuid"`
}

type ImageAmountResponse struct {
	Code  int `json:"code"`
	Count int `json:"count"`
}

type ImageQueryResponse struct {
	Code  int      `json:"code"`
	UUIDs []string `json:"uuid"`
}

type UserResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

type SystemRegister struct {
	Register string `json:"register"`
}

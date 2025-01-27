package handlers

type ImageResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type ImageQueryResponse struct {
	Code  int      `json:"code"`
	UUIDs []string `json:"uuid"`
}

type UserResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

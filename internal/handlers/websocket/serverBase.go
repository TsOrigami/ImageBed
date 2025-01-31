package websocket

import (
	dbUser "ImageV2/internal/db/user"
	Opt "ImageV2/internal/handlers/operate"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

// websocket 升级并跨域
var (
	upgrade = &websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// WebSocketBase TODO:服务基本函数
func WebSocketBase(w http.ResponseWriter, r *http.Request) {
	var (
		err         error
		conn        *websocket.Conn
		ws          *WsConn
		loginInfo   *Opt.LoginResponse
		dataOperate map[string][]string
		username    string
		token       string
	)
	if dataOperate, err = getLoginData(r); err != nil {
		fmt.Printf("err: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	loginInfo, err = Opt.LoginOperate(dataOperate)
	if conn, err = upgrade.Upgrade(w, r, nil); err != nil {
		return
	}
	token = loginInfo.ResData.Token
	if username, err = dbUser.GetUsername(token); err == nil {
		dataOperate["username"] = append(dataOperate["username"], username)
	}
	if ws, err = InitWebSocket(conn, username); err != nil {
		return
	}
	// 使得inChan和outChan耦合起来
	for {
		var data []byte
		if data, err = ws.InChanRead(); err != nil {
			goto ERR
		}
		if err = ws.OutChanWrite(data); err != nil {
			goto ERR
		}
	}
ERR:
	ws.CloseConn()
}

func getLoginData(r *http.Request) (map[string][]string, error) {
	dataOperate := make(map[string][]string)
	if r.Method != http.MethodGet {
		fmt.Printf(r.Method)
		return nil, fmt.Errorf("method Not Allowed")
	}
	queryParams := r.URL.Query()
	for key, values := range queryParams {
		dataOperate[key] = append(dataOperate[key], values...)
	}
	for key, values := range r.Header {
		dataOperate[key] = append(dataOperate[key], values...)
	}
	return dataOperate, nil
}

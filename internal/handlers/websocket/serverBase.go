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

// WsServerBase 服务基本函数
func WsServerBase(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("1. 收到WebSocket连接请求: %s\n", r.URL.String())
	fmt.Printf("2. 请求头信息: %+v\n", r.Header)

	var (
		err         error
		conn        *websocket.Conn
		ws          *WsConn
		loginInfo   *Opt.LoginResponse
		dataOperate map[string][]string
		username    string
		token       string
	)

	// 检查upgrade是否正确初始化
	if upgrade == nil {
		fmt.Println("错误: upgrade未初始化")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 首先进行 WebSocket 握手
	fmt.Println("3. 尝试进行WebSocket握手...")
	if conn, err = upgrade.Upgrade(w, r, nil); err != nil {
		fmt.Printf("4. WebSocket握手失败: %v\n", err)
		http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
		return
	}
	fmt.Println("4. WebSocket握手成功")

	// 获取登录数据
	fmt.Println("5. 尝试获取登录数据...")
	if dataOperate, err = getLoginData(r); err != nil {
		fmt.Printf("6. 获取登录数据失败: %v\n", err)
		err := conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error()))
		if err != nil {
			return
		}
		err = conn.Close()
		if err != nil {
			return
		}
		return
	}
	fmt.Printf("6. 获取到的登录数据: %+v\n", dataOperate)

	// 验证登录信息
	fmt.Println("7. 尝试验证登录信息...")
	if loginInfo, err = Opt.LoginOperate(dataOperate); err != nil {
		fmt.Printf("8. 登录验证失败: %v\n", err)
		err := conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error()))
		if err != nil {
			return
		}
		err = conn.Close()
		if err != nil {
			return
		}
		return
	}
	fmt.Printf("8. 登录验证成功: %+v\n", loginInfo)

	token = loginInfo.ResData.Token
	fmt.Printf("9. 获取到的token: %s\n", token)

	if username, err = dbUser.GetUsername(token); err != nil {
		fmt.Printf("10. 获取用户名失败: %v\n", err)
	} else {
		fmt.Printf("10. 获取到的用户名: %s\n", username)
		dataOperate["username"] = append(dataOperate["username"], username)
	}

	fmt.Println("11. 尝试初始化WebSocket连接...")
	if ws, err = InitWebSocket(conn, username); err != nil {
		fmt.Printf("12. 初始化WebSocket失败: %v\n", err)
		err := conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Error()))
		if err != nil {
			return
		}
		err = conn.Close()
		if err != nil {
			return
		}
		return
	}
	fmt.Println("12. WebSocket初始化成功")

	// 使得inChan和outChan耦合起来
	fmt.Println("13. 开始处理WebSocket消息...")
	for {
		var data []byte
		if data, err = ws.InChanRead(); err != nil {
			fmt.Printf("读取消息失败: %v\n", err)
			goto ERR
		}
		if err = ws.OutChanWrite(data); err != nil {
			fmt.Printf("写入消息失败: %v\n", err)
			goto ERR
		}
	}
ERR:
	fmt.Println("14. WebSocket连接关闭")
	ws.CloseConn()
}

func getLoginData(r *http.Request) (map[string][]string, error) {
	dataOperate := make(map[string][]string)
	queryParams := r.URL.Query()
	for key, values := range queryParams {
		dataOperate[key] = append(dataOperate[key], values...)
	}
	for key, values := range r.Header {
		dataOperate[key] = append(dataOperate[key], values...)
	}
	return dataOperate, nil
}

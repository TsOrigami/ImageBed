package websocket

import (
	Opt "ImageV2/internal/handlers/operate"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

// WsConn TODO:封装的基本结构体
type WsConn struct {
	inChan    chan []byte
	outChan   chan []byte
	closeChan chan []byte
	isClose   bool // 通道closeChan是否已经关闭
	username  string
	mutex     sync.Mutex
	conn      *websocket.Conn
}

// HandelWebSocket 注册WebSocket路由
func HandelWebSocket() {
	http.HandleFunc("/ws", WsServerBase)
}

// InitWebSocket TODO:初始化Websocket
func InitWebSocket(conn *websocket.Conn, user string) (ws *WsConn, err error) {
	ws = &WsConn{
		inChan:    make(chan []byte, 1024),
		outChan:   make(chan []byte, 1024),
		closeChan: make(chan []byte, 1024),
		username:  user,
		conn:      conn,
	}
	// 完善必要协程：读取客户端数据协程/发送数据协程
	go ws.readMsgLoop()
	go ws.writeMsgLoop()
	return
}

// InChanRead TODO:读取inChan的数据
func (conn *WsConn) InChanRead() (data []byte, err error) {
	select {
	case data = <-conn.inChan:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

// InChanWrite TODO:inChan写入数据
func (conn *WsConn) InChanWrite(data []byte) (err error) {
	select {
	case conn.inChan <- data:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

// OutChanRead TODO:读取inChan的数据
func (conn *WsConn) OutChanRead() (data []byte, err error) {
	select {
	case data = <-conn.outChan:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

// OutChanWrite TODO:inChan写入数据
func (conn *WsConn) OutChanWrite(data []byte) (err error) {
	select {
	case conn.outChan <- data:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

// CloseConn TODO:关闭WebSocket连接 仅此一次
func (conn *WsConn) CloseConn() {
	// 关闭closeChan以控制inChan/outChan策略,仅此一次
	conn.mutex.Lock()
	if !conn.isClose {
		close(conn.closeChan)
		conn.isClose = true
	}
	conn.mutex.Unlock()
	//关闭WebSocket的连接,conn.Close()是并发安全可以多次关闭
	_ = conn.conn.Close()
}

// readMsgLoop TODO:读取客户端发送的数据写入到inChan
func (conn *WsConn) readMsgLoop() {
	for {
		// 确定数据结构
		var (
			data []byte
			err  error
		)
		// 接受数据
		if _, data, err = conn.conn.ReadMessage(); err != nil {
			goto ERR
		}
		// 写入数据
		if err = conn.InChanWrite(data); err != nil {
			goto ERR
		}
	}
ERR:
	conn.CloseConn()
}

// writeMsgLoop TODO:读取outChan的数据响应给客户端
func (conn *WsConn) writeMsgLoop() {
	for {
		var (
			data    []byte
			err     error
			resData *WsResponse
		)
		// 读取数据
		if data, err = conn.OutChanRead(); err != nil {
			goto ERR
		}
		// 解析数据
		if resData, err = operateData(data, conn); err != nil {
			goto ERR
		}
		switch resData.method {
		case "upload":
			data, err = json.Marshal(resData.ResDataUpload)
		case "delete":
			data, err = json.Marshal(resData.ResDataDelete)
		case "query":
			data, err = json.Marshal(resData.ResDataQuery)
		default:
			err = errors.New("method not found")
		}
		// 发送数据
		if err = conn.conn.WriteMessage(1, data); err != nil {
			goto ERR
		}
	}
ERR:
	conn.CloseConn()
}

type WsResponse struct {
	ResDataUpload *Opt.UploadResponse `json:"ResData_upload"`
	ResDataDelete *Opt.DeleteResponse `json:"ResData_delete"`
	ResDataQuery  *Opt.QueryResponse  `json:"ResData_query"`
	method        string
}

func operateData(data []byte, conn *WsConn) (*WsResponse, error) {
	var (
		username string
		err      error
		method   string
	)
	dataJson := make(map[string]string)
	dataOperate := make(map[string][]string)
	username = conn.username
	if err = json.Unmarshal(data, &dataJson); err != nil {
		return nil, err
	}
	method = dataJson["method"]
	for key, value := range dataJson {
		dataOperate[key] = append(dataOperate[key], value)
	}
	dataOperate["username"] = append(dataOperate["username"], username)
	fmt.Printf("dataOperate: %v\n", dataOperate)
	response := &WsResponse{}
	UploadRes := &Opt.UploadResponse{}
	DeleteRes := &Opt.DeleteResponse{}
	QueryRes := &Opt.QueryResponse{}
	switch method {
	case "upload":
		UploadRes, err = Opt.UploadOperate(dataOperate)
	case "delete":
		DeleteRes, err = Opt.DeleteOperate(dataOperate)
	case "query":
		QueryRes, err = Opt.QueryOperate(dataOperate)
	default:
		err = errors.New("method not found")
	}
	response = &WsResponse{
		ResDataUpload: UploadRes,
		ResDataDelete: DeleteRes,
		ResDataQuery:  QueryRes,
		method:        method,
	}
	return response, err
}

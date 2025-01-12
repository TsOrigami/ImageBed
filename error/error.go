package error

// 错误处理

import (
	"fmt"
	"log"
	"net/http"
)

// LogError 记录错误并返回自定义错误信息
func LogError(err error, msg string) {
	if err != nil {
		log.Printf("[ERROR] %s: %v", msg, err)
	}
}

// RespondWithError 返回标准化的 HTTP 错误响应
func RespondWithError(w http.ResponseWriter, statusCode int, errMsg string) {
	http.Error(w, errMsg, statusCode)
}

// HandleAndRespond 统一处理错误并返回响应
func HandleAndRespond(w http.ResponseWriter, err error, statusCode int, errMsg string) {
	if err != nil {
		LogError(err, errMsg)
		RespondWithError(w, statusCode, errMsg)
	}
}

func UploadError(w http.ResponseWriter) {
	RespondWithError(w, http.StatusNotFound, "上传文件失败")
}

// NotFoundError 处理404错误
func NotFoundError(w http.ResponseWriter) {
	RespondWithError(w, http.StatusNotFound, "资源未找到")
}

// BadRequestError 处理400错误
func BadRequestError(w http.ResponseWriter) {
	RespondWithError(w, http.StatusBadRequest, "无效的请求参数")
}

func MethodNotAllowedError(w http.ResponseWriter) {
	RespondWithError(w, http.StatusMethodNotAllowed, "不允许的请求方法")
}

// InternalServerError 处理500错误
func InternalServerError(w http.ResponseWriter, err error) {
	HandleAndRespond(w, err, http.StatusInternalServerError, "服务器内部错误")
}

// DatabaseError 处理数据库相关错误
func DatabaseError(w http.ResponseWriter, err error) {
	HandleAndRespond(w, err, http.StatusInternalServerError, "数据库操作失败")
}

// InvalidMD5Error 处理无效的MD5错误
func InvalidMD5Error(w http.ResponseWriter) {
	RespondWithError(w, http.StatusBadRequest, "无效的 MD5 值")
}

// InvalidParameterError 处理无效的参数错误
func InvalidParameterError(w http.ResponseWriter, paramName string) {
	RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("无效的参数: %s", paramName))
}

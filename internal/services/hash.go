package services

import (
	"crypto/sha256"
	"encoding/hex"
	uuid "github.com/satori/go.uuid"
	"io"
	"os"
)

func GetSha256(FilePath string) (string, error) {
	sha256file, err := os.Open(FilePath)
	if err != nil {
		return "", err
	}
	defer func(sha256file *os.File) {
		err := sha256file.Close()
		if err != nil {
			return
		}
	}(sha256file)
	hash := sha256.New()
	_, err = io.Copy(hash, sha256file)
	if err != nil {
		return "", err
	}
	picSha256 := hex.EncodeToString(hash.Sum(nil))
	return picSha256, nil
}

// GetUUIDv1 获取UUIDv1。v1基于时间戳和MAC地址生成
func GetUUIDv1() string {
	return uuid.NewV1().String()
}

// GetUUIDv4 获取UUIDv4。v4是随机生成的
func GetUUIDv4() string {
	return uuid.NewV4().String()
}

package services

import (
	"crypto/sha256"
	"encoding/hex"
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

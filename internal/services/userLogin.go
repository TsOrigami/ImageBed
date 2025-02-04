package services

import (
	usr "ImageV2/internal/db/user"
	"crypto/md5"
	"encoding/hex"
	"errors"
)

func Login(account string, salt string, sign string) (string, string, error) {
	userName, passwd, err := usr.GetLoginInfo(account)
	if err != nil {
		return "", "", err
	}
	hashMd5 := md5.New()
	hashMd5.Write([]byte(account + salt + passwd))
	signInquire := hex.EncodeToString(hashMd5.Sum(nil))
	if signInquire == sign {
		token := GetUUIDv1()
		err := usr.SetUserToken(account, token)
		if err != nil {
			return "", "", err
		}
		return token, userName, nil
	}
	return "", "", errors.New("登录失败")
}

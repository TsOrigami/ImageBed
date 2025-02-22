package services

import (
	sql2 "ImageV2/internal/db/sql"
	"errors"
	"time"
)

func Registered(account string, username string, passwd string) error {
	dbInfo, err := sql2.GetDB()
	if err != nil {
		return err
	}
	var userName string
	registeredTime := time.Now()
	uuidV1 := GetUUIDv1()
	_ = dbInfo.Connect.QueryRow(`SELECT user_name FROM user_info WHERE account = ?`, account).Scan(&userName)
	if userName != "" {
		return errors.New("该账号已被注册")
	}
	insertSQL := `INSERT INTO user_info (uuid, user_name, account, password, count, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err = dbInfo.Connect.Exec(insertSQL, uuidV1, username, account, passwd, 0, registeredTime)
	if err != nil {
		return err
	}
	return nil
}

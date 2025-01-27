package user

import (
	sql2 "ImageV2/internal/db/sql"
	"fmt"
)

func GetLoginInfo(account string) (string, string, error) {
	dbInfo, err := sql2.GetDB()
	if err != nil {
		fmt.Println("获取数据库连接失败: ", err)
		return "", "", err
	}
	var passwd, userName string
	err = dbInfo.Connect.QueryRow(`SELECT password, user_name FROM user_info WHERE account = ?`, account).Scan(&passwd, &userName)
	if err != nil {
		fmt.Println("查询数据失败: ", err)
		return "", "", err
	}
	return userName, passwd, nil
}

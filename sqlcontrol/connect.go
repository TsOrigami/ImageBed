package sqlcontrol

import (
	conf "ImageV2/config"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"sync"
)

var db *sql.DB
var dbErr error
var once sync.Once

// 获取MySQL连接的DSN
func GetMySqlDSN() (string, error) {
	jsonData, err := conf.GetConfigGroupAsJSON("mysql")
	if err != nil {
		return "", err
	}

	var config map[string]string
	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		return "", err
	}

	host := config["host"]
	port := config["port"]
	user := config["user"]
	password := config["password"]
	dbname := config["dbname"]

	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname
	return dsn, nil
}

// 连接数据库，如果连接失败，返回错误
func ConnectionDB() (*sql.DB, error) {
	mysqlDSN, err := GetMySqlDSN()
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接信息失败: %v", err)
	}

	once.Do(func() {
		db, dbErr = sql.Open("mysql", mysqlDSN)
		if dbErr != nil {
			fmt.Println("连接数据库失败:", dbErr)
			db = nil
		}
	})

	if db == nil {
		return nil, fmt.Errorf("数据库连接未初始化成功")
	}

	return db, nil
}

// 断开数据库连接
func DisconnectDB() {
	if db != nil {
		err := db.Close()
		if err != nil {
			fmt.Println("关闭数据库连接失败:", err)
		}
	}
}

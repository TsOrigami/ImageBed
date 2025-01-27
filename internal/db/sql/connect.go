package sql

import (
	conf "ImageV2/configs"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
)

type DB struct {
	Dsn     string
	Connect *sql.DB
}

var (
	instance *DB
	mu       sync.Mutex
)

// GetDB 获取数据库连接
func GetDB() (*DB, error) {
	if instance == nil {
		mu.Lock()
		defer mu.Unlock()
		if instance == nil {
			db, mysqlDSN, err := ConnectionDB()
			if err != nil {
				return nil, err
			}
			err = createTables(db)
			if err != nil {
				return nil, err
			}
			err = createDelTables(db)
			if err != nil {
				return nil, err
			}
			err = createUserTable(db)
			if err != nil {
				return nil, err
			}
			instance = &DB{
				Dsn:     mysqlDSN,
				Connect: db,
			}
		}
	}
	return instance, nil
}

// Close 关闭数据库连接
func (db *DB) Close() error {
	if instance != nil {
		err := DisconnectDB(db.Connect)
		if err != nil {
			return err
		}
		instance = nil
		return nil
	}
	return fmt.Errorf("数据库连接未初始化")
}

// ConnectionDB 连接数据库，如果连接失败，返回错误
func ConnectionDB() (*sql.DB, string, error) {
	var db *sql.DB
	var err error
	var jsonData []byte
	jsonData, err = conf.GetConfigGroupAsJSON("mysql")
	if err != nil {
		return nil, "", fmt.Errorf("获取数据库连接信息失败: %v", err)
	}
	var config map[string]string
	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		return nil, "", fmt.Errorf("解析数据库连接信息失败: %v", err)
	}
	host := config["host"]
	port := config["port"]
	user := config["user"]
	password := config["password"]
	dbname := config["dbname"]
	mysqlDSN := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname
	db, err = sql.Open("mysql", mysqlDSN)
	if err != nil {
		return nil, "", fmt.Errorf("连接数据库失败: %v", err)
	}
	if db == nil {
		return nil, "", fmt.Errorf("数据库连接未初始化成功")
	}
	return db, mysqlDSN, nil
}

// DisconnectDB 断开数据库连接
func DisconnectDB(db *sql.DB) error {
	if db != nil {
		err := db.Close()
		if err != nil {
			fmt.Println("关闭数据库连接失败:", err)
			return err
		}
	}
	return nil
}

// createTables 创建图片信息表
func createTables(db *sql.DB) error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS image_info (
    	uuid VARCHAR(36) PRIMARY KEY,
		image_name VARCHAR(255),
    	user_name VARCHAR(255),
		sha256Hash CHAR(64),
		created_at DATETIME
	)`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("创建表失败: %v", err)
	}
	return nil
}

// createDelTables 创建删除的图片信息表
func createDelTables(db *sql.DB) error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS image_info_del (
    	uuid VARCHAR(36) PRIMARY KEY,
		image_name VARCHAR(255),
    	user_name VARCHAR(255),
		sha256Hash CHAR(64),
		created_at DATETIME
	)`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("创建表失败: %v", err)
	}
	return nil
}

// createUserTable 创建用户信息表
func createUserTable(db *sql.DB) error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS user_info (
    	uuid VARCHAR(36) PRIMARY KEY,
		user_name VARCHAR(255),
    	account VARCHAR(255),
		password CHAR(64),
		created_at DATETIME
	)`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("创建表失败: %v", err)
	}
	return nil
}

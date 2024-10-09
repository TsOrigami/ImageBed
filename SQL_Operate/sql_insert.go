package SQL_Operate

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// ImageInfo 用于存储图片的名字、路径、MD5值和创建时间
type ImageInfo struct {
	ImageName string
	ImageURL  string
	MD5Hash   string
	CreatedAt string
}

// SaveImageInfo 保存图片信息到数据库
func SaveImageInfo(imageName, md5Hash string, createdAt time.Time) error {
	// 连接数据库
	dsn := mysqlDSN
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("无法连接到数据库: %v", err)
	}
	defer db.Close()

	// 创建表（如果不存在）
	createTableSQL := `CREATE TABLE IF NOT EXISTS image_info (
		image_name VARCHAR(255),
		md5_hash CHAR(32) PRIMARY KEY,
		created_at DATETIME
	)`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("创建表失败: %v", err)
	}

	// 插入数据
	insertSQL := `INSERT INTO image_info (image_name, md5_hash, created_at) VALUES (?, ?, ?)`
	_, err = db.Exec(insertSQL, imageName, md5Hash, createdAt)
	if err != nil {
		return fmt.Errorf("插入数据失败: %v", err)
	}

	fmt.Printf("成功插入图片信息: image_name=%s, md5_hash=%s, created_at=%s\n", imageName, md5Hash, createdAt)
	return nil
}

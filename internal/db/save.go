package db

import (
	"ImageV2/internal/services"
	"fmt"
	"time"
)

func SaveInfoToSQL(imageName, sha256Hash string, createdAt time.Time) error {
	// 获取数据库连接
	dbInfo, err := GetDB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %v", err)
	}
	// 插入数据
	uuidV1 := services.GetUUIDv1()
	insertSQL := `INSERT INTO image_info (uuid, image_name, sha256Hash, created_at) VALUES (?, ?, ?, ?)`
	_, err = dbInfo.Connect.Exec(insertSQL, uuidV1, imageName, sha256Hash, createdAt)
	if err != nil {
		fmt.Println(sha256Hash)
		return fmt.Errorf("插入数据失败: %v", err)
	}
	fmt.Printf("成功插入图片信息: image_name=%s, sha256_hash=%s, created_at=%s\n", imageName, sha256Hash, createdAt)
	return nil
}

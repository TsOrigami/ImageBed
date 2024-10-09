package SQL_Operate

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// GetImageInfoByMD5 根据 MD5 值查询图片信息
func GetImageInfoByMD5(md5Value string) (ImageInfo, error) {
	// 定义数据库连接的DSN
	dsn := "root:762005@tcp(127.0.0.1:3306)/goStudy_imgDB"

	// 连接数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return ImageInfo{}, fmt.Errorf("无法连接到数据库: %v", err)
	}
	defer db.Close()

	// 根据 MD5 值查询 image_name
	querySQL := `SELECT image_name FROM image_info WHERE md5_hash = ?`
	var imageName string
	err = db.QueryRow(querySQL, md5Value).Scan(&imageName)
	if err != nil {
		if err == sql.ErrNoRows {
			return ImageInfo{}, fmt.Errorf("找不到对应的图片")
		}
		return ImageInfo{}, fmt.Errorf("查询数据库失败: %v", err)
	}

	// 返回图片信息
	return ImageInfo{
		ImageName: imageName,
	}, nil
}

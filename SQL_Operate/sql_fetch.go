package SQL_Operate

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// FetchImageByMD5 根据MD5值获取对应的图片文件
func FetchImageByMD5(md5Hash string) (string, error) {
	// 连接数据库
	dsn := mysqlDSN
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return "", fmt.Errorf("无法连接到数据库: %v", err)
	}
	defer db.Close()

	var imageURL string
	querySQL := `SELECT image_name FROM image_info WHERE md5_hash = ?`
	err = db.QueryRow(querySQL, md5Hash).Scan(&imageURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("未找到对应的MD5值: %s", md5Hash)
		}
		return "", fmt.Errorf("查询失败: %v", err)
	}
	imageURL = ImagesPath + "/" + imageURL

	return imageURL, nil
}

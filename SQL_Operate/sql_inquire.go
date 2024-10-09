package SQL_Operate

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// GetImageInfo 查询指定范围内的图片信息（名字、路径、MD5值、创建时间）
func GetImageInfo(start, end int) ([]ImageInfo, error) {

	dsn := mysqlDSN
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("无法连接到数据库: %v", err)
	}
	defer db.Close()

	querySQL := `SELECT image_name, md5_hash, created_at FROM image_info ORDER BY md5_hash LIMIT ? OFFSET ?`
	rows, err := db.Query(querySQL, end-start+1, start-1)
	if err != nil {
		return nil, fmt.Errorf("查询数据失败: %v", err)
	}
	defer rows.Close()

	var imageInfos []ImageInfo
	for rows.Next() {
		var imageName, md5Hash string
		var createdAtStr string

		if err := rows.Scan(&imageName, &md5Hash, &createdAtStr); err != nil {
			return nil, fmt.Errorf("扫描结果失败: %v", err)
		}

		createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("时间解析失败: %v", err)
		}

		formattedCreatedAt := createdAt.Format("2006-01-02 15:04:05")

		imageURL := ImagesPath + "/" + imageName

		imageInfos = append(imageInfos, ImageInfo{
			ImageName: imageName,
			ImageURL:  imageURL,
			MD5Hash:   md5Hash,
			CreatedAt: formattedCreatedAt, // 返回格式化后的时间
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历结果失败: %v", err)
	}

	return imageInfos, nil
}

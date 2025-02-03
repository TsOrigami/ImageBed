package image

import (
	sql2 "ImageV2/internal/db/sql"
	"database/sql"
	"fmt"
)

type PicInfo struct {
	UUID       string
	ImageName  string
	Username   string
	Sha256Hash string
	CreatedAt  string
}

func GetInfoByUUID(uuidInquire string) (PicInfo, error) {
	imageInfo := PicInfo{}
	dbInfo, err := sql2.GetDB()
	if err != nil {
		return imageInfo, fmt.Errorf("获取数据库连接失败: %v", err)
	}
	rows, err := dbInfo.Connect.Query("SELECT * FROM image_info WHERE uuid = ?", uuidInquire)
	if err != nil {
		return imageInfo, fmt.Errorf("查询数据失败: %v", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)
	for rows.Next() {
		err := rows.Scan(&imageInfo.UUID, &imageInfo.ImageName, &imageInfo.Username, &imageInfo.Sha256Hash, &imageInfo.CreatedAt)
		if err != nil {
			return imageInfo, fmt.Errorf("获取数据失败: %v", err)
		}
	}
	return imageInfo, nil
}

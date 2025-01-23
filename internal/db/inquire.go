package db

import (
	"database/sql"
	"fmt"
)

type PicInfo struct {
	UUID       string
	ImageName  string
	Sha256Hash string
	CreatedAt  string
}

func GetInfoByUUID(uuidInquire string) (PicInfo, error) {
	imageInfo := PicInfo{}
	dbInfo, err := GetDB()
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
	var isDelect = true
	for rows.Next() {
		err := rows.Scan(&imageInfo.UUID, &imageInfo.ImageName, &imageInfo.Sha256Hash, &imageInfo.CreatedAt, &isDelect)
		if err != nil {
			return imageInfo, fmt.Errorf("获取数据失败: %v", err)
		}
		if !isDelect {
			break
		}
	}
	if isDelect {
		return imageInfo, fmt.Errorf("未查询到数据")
	} else {
		return imageInfo, nil
	}

}

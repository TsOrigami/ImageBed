package image

import (
	sql2 "ImageV2/internal/db/sql"
	"database/sql"
	"errors"
	"fmt"
)

func DeleteInfoFromSQL(uuidDel string) error {
	// 获取数据库连接
	dbInfo, err := sql2.GetDB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %v", err)
	}

	// 查询要删除的记录信息
	var uuid, imageName, sha256Hash string
	var createdAt []uint8
	// 查询对应的记录 于 image_info 表
	querySQL := `SELECT uuid, image_name, sha256Hash, created_at FROM image_info WHERE uuid = ?`
	err = dbInfo.Connect.QueryRow(querySQL, uuidDel).Scan(&uuid, &imageName, &sha256Hash, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("未找到对应的记录，UUID: %s", uuidDel)
		}
		return fmt.Errorf("查询数据失败: %v", err)
	}

	// 将记录插入 到 image_info_del 表
	insertSQL := `INSERT INTO image_info_del (uuid, image_name, sha256Hash, created_at) VALUES (?, ?, ?, ?)`
	_, err = dbInfo.Connect.Exec(insertSQL, uuid, imageName, sha256Hash, createdAt)
	if err != nil {
		return fmt.Errorf("将数据插入到删除表失败: %v", err)
	}

	// 从 image_info 表中删除记录
	deleteSQL := `DELETE FROM image_info WHERE uuid = ?`
	_, err = dbInfo.Connect.Exec(deleteSQL, uuidDel)
	if err != nil {
		return fmt.Errorf("删除数据失败: %v", err)
	}

	fmt.Printf("成功删除图片信息，UUID: %s\n", uuidDel)
	return nil
}

package db

import (
	"fmt"
)

func DeleteInfoFromSQL(uuidDel string) error {
	// 获取数据库连接
	dbInfo, err := GetDB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %v", err)
	}
	// 插入数据
	updateSQL := `UPDATE image_info SET is_deleted = true WHERE uuid = ?`
	_, err = dbInfo.Connect.Exec(updateSQL, uuidDel)
	if err != nil {
		return fmt.Errorf("删除数据失败: %v", err)
	}
	fmt.Printf("成功删除图片信息")
	return nil
}

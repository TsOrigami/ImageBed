package SQL_Operate

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

//该函数接受一个参数：MD5值，类型为字符串
//该函数通过MD5值删除数据库中的记录，并删除本地文件
//如果成功删除，则返回nil；否则返回错误

func DeleteImageByMD5(md5Hash string) error {
	// 连接数据库
	dsn := mysqlDSN
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("无法连接到数据库: %v", err)
	}
	defer db.Close()

	// 查询对应的 image_name
	var imageURL string
	querySQL := `SELECT image_name FROM image_info WHERE md5_hash = ?`
	err = db.QueryRow(querySQL, md5Hash).Scan(&imageURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("未找到对应的MD5值: %s", md5Hash)
		}
		return fmt.Errorf("查询失败: %v", err)
	}

	// 构造完整的文件路径
	fullPath := ImagesPath + "/" + imageURL

	// 删除本地文件
	err = os.Remove(fullPath)
	if err != nil {
		return fmt.Errorf("删除文件失败: %v", err)
	}

	// 删除数据库中的记录
	deleteSQL := `DELETE FROM image_info WHERE md5_hash = ?`
	_, err = db.Exec(deleteSQL, md5Hash)
	if err != nil {
		return fmt.Errorf("删除数据库记录失败: %v", err)
	}

	fmt.Printf("成功删除MD5为 %s 的图片和数据库记录\n", md5Hash)
	return nil
}

package db

import (
	"fmt"
	"strconv"
)

func GetInfoQuery(startStr string, endStr string) ([]string, error) {
	var err error
	db, err := GetDB()
	if err != nil {
		return nil, fmt.Errorf("获取数据库连接失败: %v", err)
	}
	start, err := strconv.Atoi(startStr)
	end, err := strconv.Atoi(endStr)
	if err != nil {
		return nil, fmt.Errorf("start 或 end 参数错误: %v", err)
	}
	if start > end {
		return nil, fmt.Errorf("start 不能大于 end")
	}
	limit := end - start + 1
	offset := start - 1
	rows, err := db.Connect.Query("SELECT uuid FROM image_info LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("查询数据失败: %v", err)
	}
	var uuids []string
	for rows.Next() {
		var uuid string
		if err := rows.Scan(&uuid); err != nil {
			return nil, fmt.Errorf("failed to read data: %v", err)
		}
		uuids = append(uuids, uuid)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %v", err)
	}
	return uuids, nil
}

package services

import (
	conf "ImageV2/configs"
	"encoding/json"
	"fmt"
)

func GetSavePath() (string, error) {
	var err error
	var jsonData []byte
	jsonData, err = conf.GetConfigGroupAsJSON("server")
	if err != nil {
		return "", fmt.Errorf("获取配置信息失败: %v", err)
	}
	var config map[string]string
	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		return "", fmt.Errorf("解析配置信息失败: %v", err)
	}
	return config["imagesPath"], nil
}

package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Config map[string]string

func GetConfigGroupAsJSON(groupName string) ([]byte, error) {
	file, err := os.Open("./config/config.conf")
	if err != nil {
		return nil, fmt.Errorf("无法打开 ./config/config.conf 文件: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	var groupConfig = make(Config)
	inGroup := false
	hasGroup := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if len(line) == 0 || line[0] == '#' {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			group := strings.Trim(line, "[]")
			if group == groupName {
				inGroup = true
				hasGroup = true
			} else {
				inGroup = false
			}
			continue
		}

		if inGroup {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
					value = value[1 : len(value)-1]
				}
				groupConfig[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if !hasGroup {
		return nil, errors.New("未找到指定的配置组")
	}

	jsonData, err := json.Marshal(groupConfig)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

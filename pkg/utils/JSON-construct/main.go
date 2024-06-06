package main

import (
	"encoding/json"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	argMap := make(map[string]interface{})

	for _, arg := range args {
		// 分割参数名和参数值
		splitArg := strings.SplitN(arg, "=", 2)
		if len(splitArg) != 2 {
			continue
		}
		argName := splitArg[0]
		argValue := splitArg[1]

		// 尝试读取文件
		data, err := os.ReadFile(argValue)
		if err == nil {
			// 如果成功，将文件内容转换为字节数组
			argMap[argName] = data
		} else {
			// 否则，直接使用参数值
			argMap[argName] = argValue
		}
	}

	// 将结果转换为 JSON
	jsonData, err := json.Marshal(argMap)
	if err != nil {
		panic(err)
	}

	// 将 JSON 数据写入文件
	err = os.WriteFile("args.json", jsonData, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

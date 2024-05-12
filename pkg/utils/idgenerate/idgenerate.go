package idgenerate

import (
	"github.com/google/uuid"
)

const (
	Prefix          = "MINK8S-"
	NamespacePrefix = Prefix + "NAMESPACE-"
	PodPrefix       = Prefix + "POD-"
)

func GenerateID() (string, error) {
	// rand.NewSource(time.Now().UnixNano())

	// // 定义字符池，包含数字和字母
	// const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// // 生成长度为8的随机ID
	// var id strings.Builder
	// for i := 0; i < 8; i++ {
	// 	// 随机选择字符池中的一个字符
	// 	index := rand.Int() % len(chars)
	// 	err := id.WriteByte(chars[index])
	// 	if err != nil {
	// 		fmt.Println("failed to write to id string")
	// 		return "", err
	// 	}
	// }

	// return id.String(), nil
	id := uuid.New().String()

	return id, nil
}

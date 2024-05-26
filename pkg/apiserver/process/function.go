package process

import (
	"encoding/json"
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/serverless/types/function"
)

func GetAllServerlessFunction() ([]string, error) {
	var names []string

	pairs, err := EtcdCli.GetWithPrefix(functionPrefix)
	if err != nil {
		fmt.Println(err)
		return names, err
	}

	for _, v := range pairs {
		var tmp function.Function
		err := json.Unmarshal([]byte(v.Value), &tmp)
		if err != nil {
			fmt.Println("failed to unmarshal")
			continue
		}

		names = append(names, tmp.Metadata.Name)
	}

	return names, nil
}

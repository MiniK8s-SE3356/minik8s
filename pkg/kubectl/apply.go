package cmdline

import (
	"os"

	"github.com/spf13/cobra"
)

func ApplyCmdHandler(cmd *cobra.Command, args []string) {
	result := checkFilePath(args)
	if !result {
		return
	}
}

func checkFilePath(args []string) bool {
	// 检查参数给出的文件路径是否存在

	if len(args) == 0 {
		return false
	}

	result, err := os.Stat(args[0])
	if err != nil {
		return false
	}

	if result.IsDir() {
		return false
	}

	return true
}

package cmdline

import (
	"fmt"

	"github.com/spf13/cobra"
)

const defaultNamespace = "Default"

var kubectlCmd = &cobra.Command{
	Use:   "kubectl",
	Short: "kubectl",
	Long:  `kubectl`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "apply",
	Long:  `apply`,
	Run:   ApplyCmdHandler,
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get",
	Long:  `get`,
	Run:   GetCmdHandler,
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete",
	Long:  `delete`,
	Run:   DeleteCmdHandler,
}

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "describe",
	Long:  "describe",
	Run:   DescribeCmdHandler,
}

func init() {
	kubectlCmd.AddCommand(applyCmd)
	kubectlCmd.AddCommand(getCmd)
	kubectlCmd.AddCommand(deleteCmd)
	kubectlCmd.AddCommand(describeCmd)
}

func Exec() {
	if err := kubectlCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

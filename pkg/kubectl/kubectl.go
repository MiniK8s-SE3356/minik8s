package kubectl

import (
	"fmt"

	"github.com/MiniK8s-SE3356/minik8s/pkg/kubectl/cmdline"
	"github.com/spf13/cobra"
)

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
	Run:   cmdline.ApplyCmdHandler,
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get",
	Long:  `get`,
	Run:   cmdline.GetCmdHandler,
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete",
	Long:  `delete`,
	Run:   cmdline.DeleteCmdHandler,
}

func init() {
	kubectlCmd.AddCommand(applyCmd)
	kubectlCmd.AddCommand(getCmd)
	kubectlCmd.AddCommand(deleteCmd)
}

func Exec() {
	if err := kubectlCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

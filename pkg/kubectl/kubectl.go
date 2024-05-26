package kubectl

import (
	"flag"
	"fmt"
	"time"

	"github.com/MiniK8s-SE3356/minik8s/pkg/kubectl/cmdline"
	minik8s_message "github.com/MiniK8s-SE3356/minik8s/pkg/utils/message"
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
	apiServerIP := flag.String("apiserverip", "127.0.0.1", "APIServer IP address")
	mqConfig := minik8s_message.MQConfig{
		User:       "guest",
		Password:   "guest",
		Host:       *apiServerIP,
		Port:       "5672",
		Vhost:      "/",
		MaxRetry:   5,
		RetryDelay: 5 * time.Second,
	}
	tmp, err := minik8s_message.NewMQConnection(&mqConfig)
	if err != nil {
		fmt.Println("failed to connect to mq", err)
		return
	}
	cmdline.MqConn = tmp

	if err := kubectlCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

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

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create",
	Long:  `create`,
	Run:   cmdline.CreateCmdHandler,
}

var triggerCmd = &cobra.Command{
	Use:   "trigger",
	Short: "trigger",
	Long:  `trigger`,
	Run:   cmdline.TriggerCmdHandler,
}

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "submit",
	Long:  `submit`,
	Run:   cmdline.SubmitCmdHandler,
}

func init() {
	kubectlCmd.AddCommand(applyCmd)
	kubectlCmd.AddCommand(getCmd)
	kubectlCmd.AddCommand(deleteCmd)
	kubectlCmd.AddCommand(createCmd)
	kubectlCmd.AddCommand(triggerCmd)
	kubectlCmd.AddCommand(submitCmd)
}

func Exec() {
	apiServerIP := flag.String("apiserverip", "127.0.0.1", "APIServer IP address")
	apiServerPort := flag.String("apiserverport", "8080", "APIServer port")
	cmdline.RootURL = "http://" + *apiServerIP + ":" + *apiServerPort

	serverlessIP := flag.String("serverlessip", "127.0.0.1", "serverlessip")
	serverlessPort := flag.String("serverlessport", "8081", "serverlessip")
	cmdline.ServerlessRootURL = "http://" + *serverlessIP + ":" + *serverlessPort

	gpuctlIP := flag.String("gpuctlip", "127.0.0.1", "gpuctlip")
	gpuctlPort := flag.String("gpuctlport", "8083", "gpuctlport")
	cmdline.GPUCtlRootURL = "http://" + *gpuctlIP + ":" + *gpuctlPort

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

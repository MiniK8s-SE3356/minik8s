package ssh_test

import (
	"testing"

	"github.com/MiniK8s-SE3356/minik8s/pkg/gpu/utils/ssh"
)

func TestMain(m *testing.M) {
	sshclient, err := ssh.NewSSHClient(
		ssh.SSHConfig{
			Username: "stu1045",
			Password: "RFVrGtwdKSGl",
			Hostname: "pilogin.hpc.sjtu.edu.cn",
		},
	)
	if err != nil {
		panic(err)
	}

	defer sshclient.Close()

	//...Test for batch command...//
	// cmds := []string{
	// 	"cd",
	// 	"ls -a -l",
	// 	"mkdir -p minik8s-test",
	// 	"cd minik8s-test",
	// 	"pwd",
	// 	"echo 'Hello, World!' > hello.txt",
	// 	"cat hello.txt",
	// }

	// out, err := sshclient.BatchCmd(cmds)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(out)

	//...Test for post directory...//

	// cmds := []string{
	// 	"mkdir -p minik8s-dir-test",
	// }
	// _, err = sshclient.BatchCmd(cmds)
	// if err != nil {
	// 	panic(err)
	// }

	err = sshclient.PostDirectory("/home/xubbbb/Code/CloudOS/minik8s/scripts", "minik8s-dir-test")
	if err != nil {
		panic(err)
	}
}

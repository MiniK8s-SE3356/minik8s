package ssh_test

import (
	"fmt"
	"testing"

	"github.com/MiniK8s-SE3356/minik8s/pkg/gpu/utils/ssh"
)

func TestMain(m *testing.M) {
	sshclient, err := ssh.NewSSHClient(
		ssh.SSHConfig{
			Username: "stu1045",
			Password: "RFVrGtwdKSGl",
			Hostname: "sylogin.hpc.sjtu.edu.cn",
		},
	)
	if err != nil {
		panic(err)
	}

	defer sshclient.Close()

	cmds := []string{
		"cd",
		"ls -a -l",
		"mkdir -p minik8s-test",
		"cd minik8s-test",
		"pwd",
		"echo 'Hello, World!' > hello.txt",
		"cat hello.txt",
	}

	out, err := sshclient.BatchCmd(cmds)
	if err != nil {
		panic(err)
	}

	fmt.Println(out)
}

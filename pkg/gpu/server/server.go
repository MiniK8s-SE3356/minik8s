package server

import "github.com/MiniK8s-SE3356/minik8s/pkg/gpu/utils/ssh"

type JobServer struct {
	SSHClient *ssh.SSHClient
}

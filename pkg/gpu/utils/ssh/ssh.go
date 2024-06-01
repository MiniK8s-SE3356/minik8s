package ssh

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/melbahja/goph"
)

type SSHConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Hostname string `json:"hostname"`
}

type SSHClient struct {
	Config SSHConfig
	Client *goph.Client
}

func NewSSHClient(config SSHConfig) (*SSHClient, error) {
	auth := goph.Password(config.Password)
	client, err := goph.NewUnknown(
		config.Username,
		config.Hostname,
		auth,
	)

	if err != nil {
		fmt.Println("Failed to create SSH client, err: ", err)
		return nil, err
	}

	return &SSHClient{
		Config: config,
		Client: client,
	}, nil
}

func (sc *SSHClient) BatchCmd(cmds []string) (string, error) {
	cmd := ""
	for _, c := range cmds {
		cmd += c + ";"
	}
	out, err := sc.Client.Run(cmd)
	if err != nil {
		fmt.Println("Failed to run command BATCHCMD, cmd: ", cmd)
		return "", err
	}
	fmt.Println("Command: ", cmd, " executed successfully")
	return string(out), nil
}

func (sc *SSHClient) PostDirectory(localPath, remotePath string) error {
	err := filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Failed to walk through local path")
			return err
		}

		relativePath, err := filepath.Rel(localPath, path)
		if err != nil {
			fmt.Println("Failed to get relative path")
			return err
		}
		remoteFilePath := filepath.Join(remotePath, relativePath)

		if info.IsDir() {
			out, err := sc.Client.Run("mkdir -p " + remoteFilePath)
			if err != nil {
				fmt.Println("Failed to create directory")
				return err
			}
			fmt.Println("Directory created successfully: ", string(out))
		} else {
			err := sc.Client.Upload(path, remoteFilePath)
			if err != nil {
				fmt.Println("Failed to send file")
				return err
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("Failed to send directory")
		return err
	}

	fmt.Println("Directory sent successfully")
	return nil
}

func (sc *SSHClient) PostFile(localPath, remotePath string) error {
	err := sc.Client.Upload(localPath, remotePath)
	if err != nil {
		fmt.Println("Failed to send file")
		return err
	}

	fmt.Println("File sent successfully")
	return nil
}

func (sc *SSHClient) GetFile(remotePath, localPath string) error {
	err := sc.Client.Download(remotePath, localPath)
	if err != nil {
		fmt.Println("Failed to get file")
		return err
	}

	fmt.Println("File got successfully")
	return nil
}

func (c *SSHClient) Close() error {
	return c.Client.Close()
}

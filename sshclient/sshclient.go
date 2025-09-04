package sshclient

import (
	"io"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Client wraps SSH and SFTP clients
type Client struct {
	sshClient  *ssh.Client
	sftpClient *sftp.Client
}

// NewClient creates a new SSH/SFTP client
func NewClient(host, user, password string) (*Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // For testing
	}
	sshClient, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return nil, err
	}
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		err = sshClient.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	return &Client{sshClient: sshClient, sftpClient: sftpClient}, nil
}

// Close closes SSH and SFTP connections
func (c *Client) Close() {
	err := c.sftpClient.Close()
	if err != nil {
		return
	}
	err = c.sshClient.Close()
	if err != nil {
		return
	}
}

// UploadFile uploads a file to the remote path
func (c *Client) UploadFile(remotePath string, data io.Reader) error {
	f, err := c.sftpClient.Create(remotePath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, data)
	return err
}

// ListFiles lists files in the remote directory
func (c *Client) ListFiles(dir string) ([]os.FileInfo, error) {
	return c.sftpClient.ReadDir(dir)
}

// DownloadFile downloads a file to the local path
func (c *Client) DownloadFile(remotePath, localPath string) error {
	remoteFile, err := c.sftpClient.Open(remotePath)
	if err != nil {
		return err
	}
	defer remoteFile.Close()
	localFile, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer localFile.Close()
	_, err = io.Copy(localFile, remoteFile)
	return err
}

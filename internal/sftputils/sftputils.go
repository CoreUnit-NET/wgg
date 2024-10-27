package sftputils

import (
	"errors"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SshConfig struct {
	Host       string
	Port       int
	User       string
	PrivateKey string
	Password   string
	Timeout    time.Duration
}

func (config *SshConfig) VerifySshConfig() error {
	if config.Port < 1 || config.Port > 65535 {
		return errors.New("sshconfig: invalid port number")
	}

	if len(config.Password) <= 0 && len(config.PrivateKey) <= 0 {
		return errors.New("sshconfig: password or private key is required")
	}

	if len(config.User) <= 0 {
		return errors.New("sshconfig: user is required")
	}

	if len(config.Host) <= 0 {
		return errors.New("sshconfig: host is required")
	}

	if config.Timeout <= 0 {
		return errors.New("sshconfig: timeout lesser then 1 is not allowed")
	}

	return nil
}

func (config *SshConfig) FillSshConfig() {
	if config.Port == 0 {
		config.Port = 22
	}

	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	if len(config.User) == 0 {
		config.User = "root"
	}
}

func HandleSftp(
	sshConfig *SshConfig,
	handle func(
		*sftp.Client,
		*ssh.Session,
	) error,
) error {
	// var hostkeyCallback ssh.HostKeyCallback
	// hostkeyCallback, err = knownhosts.New(homeDir + "/.ssh/known_hosts")
	// if err != nil {
	// 	return errors.New("error parsing known hosts: " + err.Error())
	// }

	sshConfig.FillSshConfig()
	err := sshConfig.VerifySshConfig()
	if err != nil {
		return errors.New("error verifying ssh config: " + err.Error())
	}

	authMethods := []ssh.AuthMethod{}

	if len(sshConfig.Password) > 0 {
		authMethods = append(authMethods, ssh.Password(sshConfig.Password))
	}

	if len(sshConfig.PrivateKey) > 0 {
		signer, err := ssh.ParsePrivateKey([]byte(sshConfig.PrivateKey))
		if err != nil {
			return errors.New("error parsing private key: " + err.Error())
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	conf := &ssh.ClientConfig{
		User: sshConfig.User,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Auth: authMethods,
	}

	// sftp
	sftpSshClient, err := ssh.Dial("tcp", sshConfig.Host+":"+strconv.Itoa(sshConfig.Port), conf)
	if err != nil {
		return errors.New("error dialing: " + err.Error())
	}
	defer sftpSshClient.Close()

	sftp, err := sftp.NewClient(
		sftpSshClient,
	)
	if err != nil {
		return errors.New("error creating sftp client: " + err.Error())
	}
	defer sftp.Close()

	// // session
	// sessionSshClient, err := ssh.Dial("tcp", sshConfig.Host+":"+strconv.Itoa(sshConfig.Port), conf)
	// if err != nil {
	// 	return errors.New("error dialing: " + err.Error())
	// }
	// defer sessionSshClient.Close()

	// session, err := sessionSshClient.NewSession()
	// if err != nil {
	// 	return errors.New("error creating ssh session: " + err.Error())
	// }
	// defer session.Close()

	// handle
	err = handle(sftp, nil)
	if err != nil {
		return errors.New("error handling: " + err.Error())
	}

	return nil
}

func JoinPath(sftp *sftp.Client, path ...string) (string, error) {
	if len(path) != 0 &&
		len(path[0]) != 0 {
		if strings.HasPrefix(path[0], "~/") ||
			strings.HasPrefix(path[0], "./") {
			cwd, err := sftp.Getwd()
			if err != nil {
				return "", errors.New("error getting cwd: " + err.Error())
			}

			if strings.HasPrefix(path[0], "../") {
				path[0] = cwd + "/" + path[0]
			} else {
				path[0] = cwd + path[0][1:]
			}
		}
	}

	return sftp.Join(path...), nil
}

package knownhost

import (
	"bufio"
	"context"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type KnownHost struct {
	knownHostFile string
}

func NewKnownHost(opts ...Option) *KnownHost {
	knownHost := &KnownHost{}
	for _, opt := range opts {
		opt(knownHost)
	}
	return knownHost
}

type Option func(host *KnownHost)

func WithDefaultKnownHostsFile(yes bool) Option {
	return func(host *KnownHost) {
		if yes {
			host.knownHostFile, _ = host.GetDefaultKnownHostFile()
		}
	}
}

func WithCustomFile(fileName string) Option {
	return func(host *KnownHost) {
		host.knownHostFile = fileName
	}
}

// GetDefaultKnownHostFile returns default knownhosts file path
func (k *KnownHost) GetDefaultKnownHostFile() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".ssh", "known_hosts"), nil
}

// ReadLocalHostKeyForHost read known host key for specify host
func (k *KnownHost) ReadLocalHostKeyForHost(host string) (hostKey ssh.PublicKey, err error) {
	file, err := os.Open(k.knownHostFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				return nil, err
			}
		}
	}
	return hostKey, nil
}

func (k *KnownHost) GetKeysForHost(host string, timeout time.Duration) ([]ssh.PublicKey, error) {
	var (
		publicKeys = make([]ssh.PublicKey, 0, len(supportedHostKeyAlgorithms))
		recv       = make(chan ssh.PublicKey, 1)
	)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ctx := context.WithValue(ctxTimeout, timeoutKey, timeout)
	ctx = context.WithValue(ctx, hostKey, host)
	for _, algorithm := range supportedHostKeyAlgorithms {
		ctx = context.WithValue(ctx, algorithmKey, algorithm)
		go processFilterData(ctx, recv)
	}

	for {
		select {
		case <-ctxTimeout.Done():
			return publicKeys, nil
		case pubKey := <-recv:
			publicKeys = append(publicKeys, pubKey)
			continue
		}
	}

	return publicKeys, nil

}

func processFilterData(ctx context.Context, recv chan ssh.PublicKey) {
	publicKey := getPublicKey(ctx)
	if publicKey != nil {
		recv <- publicKey
	}
}

func getPublicKey(ctx context.Context) (key ssh.PublicKey) {
	timeout := getTimeoutFromContext(ctx)
	host := getHostFromContext(ctx)
	algorithm := getAlgorithmFromContext(ctx)

	d := net.Dialer{Timeout: timeout}
	conn, err := d.Dial("tcp", host)
	if err != nil {
		return key
	}
	defer conn.Close()

	config := ssh.ClientConfig{
		HostKeyAlgorithms: []string{algorithm},
		HostKeyCallback:   hostKeyCallback(&key),
	}
	sshConn, _, _, err := ssh.NewClientConn(conn, host, &config)
	if err == nil {
		sshConn.Close()
	}
	return key

}
func hostKeyCallback(publicKey *ssh.PublicKey) func(hostname string, remote net.Addr, key ssh.PublicKey) error {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		*publicKey = key
		return nil
	}
}

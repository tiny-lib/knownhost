package knownhost

import "golang.org/x/crypto/ssh"

var (
	supportedHostKeyAlgorithms = []string{
		ssh.CertAlgoRSAv01,
		ssh.CertAlgoDSAv01,
		ssh.CertAlgoECDSA256v01,
		ssh.CertAlgoECDSA384v01,
		ssh.CertAlgoECDSA521v01,
		ssh.CertAlgoED25519v01,
		ssh.KeyAlgoECDSA256,
		ssh.KeyAlgoECDSA384,
		ssh.KeyAlgoECDSA521,
		ssh.KeyAlgoRSA,
		ssh.KeyAlgoDSA,
		ssh.KeyAlgoED25519,
	}
)

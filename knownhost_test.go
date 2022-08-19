package knownhost

import (
	"golang.org/x/crypto/ssh"
	"testing"
	"time"
)

func TestKnownHost_GetKeysForHost(t *testing.T) {
	host := NewKnownHost()
	keys, err := host.GetKeysForHost("gitee.com:22", 2*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range keys {
		payload := ssh.MarshalAuthorizedKey(key)
		t.Log(string(payload))
	}
}

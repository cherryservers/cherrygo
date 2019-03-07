package cherrygo

import (
	"crypto/rand"
	"crypto/rsa"
	"strconv"
	"testing"

	"golang.org/x/crypto/ssh"
)

func makeRandKey(t *testing.T) string {

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("error while generating test private key: %v", err)
	}

	pub, err := ssh.NewPublicKey(&priv.PublicKey)
	if err != nil {
		t.Fatalf("error while generating test public key: %v", err)
	}

	return string(ssh.MarshalAuthorizedKey(pub))
}

func createSSHKey(t *testing.T, c *Client) SSHKeys {

	sshCreateRequest := CreateSSHKey{
		Label: "CHERRY_TEST_LABEL-" + RandStringBytes(8),
		Key:   makeRandKey(t),
	}

	sshkey, _, err := c.SSHKey.Create(&sshCreateRequest)
	if err != nil {
		t.Fatalf("error while creating ssh key: %v", err)
	}

	return sshkey
}
func TestSSHKeyDelete(t *testing.T) {

	t.Parallel()

	c := setupClient(t)

	sshKey := createSSHKey(t, c)

	sshKeyString := strconv.Itoa(sshKey.ID)

	sshDeleteRequest := DeleteSSHKey{ID: sshKeyString}

	c.SSHKey.Delete(&sshDeleteRequest)

	sshkey, _, err := c.SSHKey.List(sshKeyString)
	if err == nil {
		t.Fatalf("it seems key wasn't deleted: %v", sshkey)
	}
}

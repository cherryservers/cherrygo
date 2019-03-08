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

func TestSSHKeyCreate(t *testing.T) {

	t.Parallel()

	c := setupClient(t)

	sshCreateRequest := CreateSSHKey{
		Label: "CHERRY_TEST_LABEL-" + RandStringBytes(8),
		Key:   makeRandKey(t),
	}

	sshkey, _, err := c.SSHKey.Create(&sshCreateRequest)
	if err != nil {
		t.Fatalf("error while creating ssh key: %v", err)
	}

	defer c.SSHKey.Delete(&DeleteSSHKey{ID: (strconv.Itoa(sshkey.ID))})

	if sshkey.Label != sshCreateRequest.Label {
		t.Fatalf("ssh keys label doesn't match, expected: %v, current: %v", sshCreateRequest.Label, sshkey.Label)
	}

	if sshkey.Key != sshCreateRequest.Key {
		t.Fatalf("ssh keys doesn'n match, expected: %v, current: %v", sshCreateRequest.Key, sshkey.Key)
	}

}

func TestSSHKeyUpdate(t *testing.T) {

	t.Parallel()

	c := setupClient(t)

	sshKey := createSSHKey(t, c)

	newLabel := sshKey.Label + "-new"

	sshUpateRequest := UpdateSSHKey{
		Label: newLabel,
	}

	sshKeyString := strconv.Itoa(sshKey.ID)

	key, _, err := c.SSHKey.Update(sshKeyString, &sshUpateRequest)
	if err != nil {
		t.Fatalf("error while updating ssh key: %v", err)
	}

	defer c.SSHKey.Delete(&DeleteSSHKey{ID: sshKeyString})

	if key.Label != sshUpateRequest.Label {
		t.Fatalf("expected label: %v, found: %v", sshUpateRequest.Label, key.Label)
	}

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

func TestSSHKeyList(t *testing.T) {

	t.Parallel()

	c := setupClient(t)

	sshKey := createSSHKey(t, c)

	sshKeyString := strconv.Itoa(sshKey.ID)

	defer c.SSHKey.Delete(&DeleteSSHKey{ID: sshKeyString})

	sshkey, _, err := c.SSHKey.List(sshKeyString)
	if err != nil {
		t.Fatalf("unable to get SSH key: %v", err)
	}

	if sshKey.ID != sshkey.ID {
		t.Fatalf("keys IDs doesn't match, expected: %v, current: %v", sshKey.ID, sshkey.ID)
	}

	if sshKey.Fingerprint != sshkey.Fingerprint {
		t.Fatalf("keys fingerprints doesn't match, expected: %v, current: %v", sshKey.Fingerprint, sshkey.Fingerprint)
	}

	if sshKey.Label != sshkey.Label {
		t.Fatalf("keys labels doesn't match, expected: %v, current: %v", sshKey.Label, sshkey.Label)
	}
}

func TestSSHKeysList(t *testing.T) {

	t.Parallel()

	c := setupClient(t)

	sshKey := createSSHKey(t, c)

	sshKeyString := strconv.Itoa(sshKey.ID)

	defer c.SSHKey.Delete(&DeleteSSHKey{ID: sshKeyString})

	sshkeys, _, err := c.SSHKeys.List()
	if err != nil {
		t.Fatalf("unable to get SSH key: %v", err)
	}

	for _, v := range sshkeys {
		if v.ID == sshKey.ID {
			return
		}
	}
	t.Fatalf("unable to find created key in a list: %v", sshKey.ID)

}

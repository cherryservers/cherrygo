package cherrygo

import (
	"math/rand"
	"testing"
	"time"
)

// RandStringBytes return random string
func RandStringBytes(n int) string {

	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	rand.Seed(time.Now().UnixNano())

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// setupClient sutups client
func setupClient(t *testing.T) *Client {

	c, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}

	return c
}

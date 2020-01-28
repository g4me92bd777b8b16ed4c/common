package chatenc

import (
	"crypto/rand"
	"io"

	"sync"

	"gitlab.com/aerth/x/hash/argon2id"
	"golang.org/x/crypto/nacl/secretbox"
)

type Crypt struct {
	pkey [32]byte
	mu   sync.Mutex
}

func Load(s string) *Crypt {
	c := &Crypt{}
	go c.Reload(s)
	return c
}

func (c *Crypt) Reload(s string) {
	c.mu.Lock()
	copy(c.pkey[:], argon2id.New(1, 1024, 1).Sum([]byte(s)))
	c.mu.Unlock()
}

func (c *Crypt) Hash() []byte {
	c.mu.Lock()
	defer c.mu.Unlock()
	return argon2id.NewDefault().Sum(c.pkey[:])
}

func (c *Crypt) Encrypt(b []byte) []byte {
	c.mu.Lock()
	defer c.mu.Unlock()
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		panic(err)
	}
	return secretbox.Seal(nonce[:], b, &nonce, &c.pkey)
}

func (c *Crypt) Decrypt(b []byte) []byte {
	if len(b) < 25 {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	var decryptNonce [24]byte
	copy(decryptNonce[:], b[:24])
	decrypted, ok := secretbox.Open(nil, b[24:], &decryptNonce, &c.pkey)
	if !ok {
		return nil
	}
	return decrypted
}

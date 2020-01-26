package chatenc

import (
	"crypto/rand"
	"encoding/hex"
	"io"

	"gitlab.com/aerth/x/hash/argon2id"
	"golang.org/x/crypto/nacl/secretbox"
)

type Crypt struct {
	pkey [32]byte
}

func Load(s string) *Crypt {
	var secretKeyBytes []byte
	var err error
	if len(s) != 64 {
		secretKeyBytes = argon2id.New(1, 1, 1).Sum([]byte(s))
	} else {
		secretKeyBytes, err = hex.DecodeString(s)
		if err != nil {
			panic(err)
		}
	}

	var secretKey [32]byte
	copy(secretKey[:], secretKeyBytes)

	c := &Crypt{
		pkey: secretKey,
	}
	return c
}

func (c *Crypt) Reload(s string) {
	var secretKeyBytes []byte
	secretKeyBytes = argon2id.New(1, 1, 1).Sum([]byte(s))
	copy(c.pkey[:], secretKeyBytes)

}

// Hash to see which channel you are on
func (c *Crypt) Hash() []byte {
	return argon2id.New(1, 1, 1).Sum(c.pkey[:])
}

func (c *Crypt) Encrypt(b []byte) []byte {
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
	var decryptNonce [24]byte
	copy(decryptNonce[:], b[:24])
	decrypted, ok := secretbox.Open(nil, b[24:], &decryptNonce, &c.pkey)
	if !ok {
		return nil
	}
	return decrypted
}

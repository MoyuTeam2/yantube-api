package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"hash"

	"golang.org/x/crypto/pbkdf2"
)

const DefaultSaltLen = 16

var (
	ErrInvalidHashAlgo     = errors.New("invalid hash algorithm")
	ErrInvalidPbkdf2Format = errors.New("invalid pbkdf2 string")
)

var (
	Pbkdf2Sha256  = MustPbkdf2WithAlgo("sha256", 32, 100000)
	Pbkdf2Sha512  = MustPbkdf2WithAlgo("sha512", 64, 100000)
	DefaultPbkdf2 = Pbkdf2Sha256
)

var hasherMapping = map[string]func() hash.Hash{
	"sha256": sha256.New,
	"sha512": sha512.New,
}

type Pbkdf2Impl struct {
	iter     int
	keyLen   int
	hashAlgo string
	hasher   func() hash.Hash
}

func NewPbkdf2WithAlgo(algo string, keyLen int, iter int) (*Pbkdf2Impl, error) {
	hasher, ok := hasherMapping[algo]
	if !ok {
		return nil, ErrInvalidHashAlgo
	}
	return &Pbkdf2Impl{iter, keyLen, algo, hasher}, nil
}

func MustPbkdf2WithAlgo(algo string, keyLen int, iter int) *Pbkdf2Impl {
	p, err := NewPbkdf2WithAlgo(algo, keyLen, iter)
	if err != nil {
		panic(err)
	}
	return p
}

func (p *Pbkdf2Impl) Hash(password []byte, salt []byte) []byte {
	return pbkdf2.Key(password, salt, p.iter, p.keyLen, p.hasher)
}

func (p *Pbkdf2Impl) FormatedString(password []byte, salt []byte) string {
	cipherBase64 := base64.StdEncoding.EncodeToString(p.Hash(password, salt))
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	return fmt.Sprintf(`%s$%d$%d$%s$%s`, p.hashAlgo, p.keyLen, p.iter, saltBase64, cipherBase64)
}

func MkPbkdf2String(password []byte, salt []byte, p *Pbkdf2Impl) string {
	return p.FormatedString(password, salt)
}

func ValidatePbkdf2Cipher(s string, password string) (bool, error) {
	var hashAlgo string
	var keyLen int
	var iter int
	var saltBase64, cipherBase64 string

	_, err := fmt.Sscanf(s, `%s$%d$%d$%s$%s`, &hashAlgo, &keyLen, &iter, &saltBase64, &cipherBase64)
	if err != nil {
		return false, ErrInvalidPbkdf2Format
	}

	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		return false, errors.Join(err, ErrInvalidPbkdf2Format)
	}

	cipherBytes, err := base64.StdEncoding.DecodeString(cipherBase64)
	if err != nil {
		return false, errors.Join(err, ErrInvalidPbkdf2Format)
	}

	p, err := NewPbkdf2WithAlgo(hashAlgo, keyLen, iter)
	if err != nil {
		return false, err
	}

	inputPasswordCipher := p.Hash([]byte(password), salt)
	return bytes.Equal(inputPasswordCipher, cipherBytes), nil
}

func RandomBytes(n int) ([]byte, error) {
	ret := make([]byte, n)
	if _, err := rand.Read(ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func RandomBytesWithDefaultLen() ([]byte, error) {
	return RandomBytes(DefaultSaltLen)
}

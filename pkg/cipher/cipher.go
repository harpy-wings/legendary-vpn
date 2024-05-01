package cipher

import (
	"bytes"
	gcipher "crypto/cipher"
	"crypto/rand"
	"errors"
	"log"

	"github.com/golang/snappy"
	"github.com/harpy-wings/legendary-vpn/pkg/xerror"
)

type Cipher interface {
	BlockSize() int
	Encrypt(msg []byte) ([]byte, error)
	Decrypt(bs []byte) ([]byte, error) //todo what is iv and ctext, add proper naming and comment.
}

// Default cipher interface for singleton pattern
var Default Cipher

var _ Cipher = &cipher{}

type cipher struct {
	block     gcipher.Block
	blockSize int
	config    struct {
		key []byte

		method Method
	}
}

func NewCipher(key []byte, ops ...Option) (Cipher, error) {
	c := new(cipher)
	c.setDefaults()

	for _, fn := range ops {
		err := fn(c)
		if err != nil {
			return nil, err
		}
	}

	err := c.init()
	if err != nil {
		return nil, err
	}

	return c, nil
}

// InitCipher initializes the singleton cipher
func InitCipher(key []byte, ops ...Option) error {
	c, err := NewCipher(key, ops...)
	if err != nil {
		return err
	}
	Default = c
	return nil
}

func (c *cipher) BlockSize() int {
	return c.block.BlockSize()
}

func (c *cipher) Encrypt(msg []byte) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			// todo add logger package
			log.Println(errors.Join(xerror.ErrCipherEncryptionFailure, err.(error)))
		}
	}()
	// compressing using snappy and encrypting data
	msg = append(
		msg[:c.blockSize],
		c.pkcs5Padding(snappy.Encode(nil, msg))...)

	// generating random bytes for IV
	rand.Read(msg[:c.blockSize])

	// creates encrypter using block and IV
	encrypter := gcipher.NewCBCEncrypter(c.block, msg[:c.blockSize])

	encrypter.CryptBlocks(msg[c.blockSize:], msg[c.blockSize:])
	return msg, nil
}

func (c *cipher) Decrypt(bs []byte) ([]byte, error) {
	var err error

	iv := bs[:c.blockSize]
	ctext := bs[c.blockSize:]

	defer func() {
		if err := recover(); err != nil {
			// todo add logger package
			log.Println(errors.Join(xerror.ErrCipherEncryptionFailure, err.(error)))
		}
	}()
	decrypter := gcipher.NewCBCDecrypter(c.block, iv)
	decrypter.CryptBlocks(ctext, ctext)
	ctext = c.pkcs5UnPadding(ctext)
	ctext, err = snappy.Decode(nil, ctext)
	if err != nil {
		return nil, err
	}
	return ctext, nil
}

func (c *cipher) pkcs5Padding(key []byte) []byte {
	padding := c.blockSize - len(key)%c.blockSize
	return append(key, bytes.Repeat([]byte{byte(padding)}, padding)...)
}

func (c *cipher) pkcs5UnPadding(v []byte) []byte {
	length := len(v)
	unPaddingIndex := int(v[length-1]) // see padding
	return v[:(length - unPaddingIndex)]
}

func (c *cipher) setDefaults() {
	c.config.method = MethodAES256
}

func (c *cipher) init() error {
	switch c.config.method {
	case MethodAES256:
		c.blockSize = blockSizes[MethodAES256]
	default:
		return xerror.ErrCipherMethodNotSupported
	}
	return nil
}

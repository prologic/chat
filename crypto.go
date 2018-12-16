package main

import (
	"crypto/cipher"
	"errors"
)

// Cipher represents a basic cipher implementation consisting of basic
// encryption and decryption methods as well as a way to determine the
// block size.
type Cipher interface {
	BlockSize() int
	Decrypt(dst, src []byte)
	Encrypt(dst, src []byte)
}

func checkSizeAndPad(value []byte, blocksize int) []byte {
	modulus := len(value) % blocksize
	if modulus != 0 {
		padnglen := blocksize - modulus
		for i := 0; i < padnglen; i++ {
			value = append(value, 0)
		}
	}
	return value
}

func encrypt(bcipher Cipher, value []byte) ([]byte, error) {
	value = checkSizeAndPad(value, bcipher.BlockSize())
	out := make([]byte, bcipher.BlockSize()+len(value))
	eiv := out[:bcipher.BlockSize()]
	ecbc := cipher.NewCBCEncrypter(bcipher, eiv)
	ecbc.CryptBlocks(out[bcipher.BlockSize():], value)
	return out, nil
}

func decrypt(dcipher Cipher, value []byte) ([]byte, error) {
	div := value[:dcipher.BlockSize()]
	decrypted := value[dcipher.BlockSize():]
	if len(decrypted)%dcipher.BlockSize() != 0 {
		return nil, errors.New("invalid blocksize")
	}
	dcbc := cipher.NewCBCDecrypter(dcipher, div)
	dcbc.CryptBlocks(decrypted, decrypted)
	return decrypted, nil
}

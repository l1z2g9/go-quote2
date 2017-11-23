package util

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
)

// because we are going to use TripleDES... therefore we Triple it!
const triplekey = "qkjl#5@2" + "md3g_s5Q" + "@FD&fawE" //%&@~sb9x

func PKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

func Encrypt(plaintext_ string) string {
	//plaintext := []byte("R9510752X8136") // Hello World! = 12 bytes.
	plaintext := []byte(plaintext_)

	block, _ := des.NewTripleDESCipher([]byte(triplekey))

	ciphertext := []byte("abcdef1234567890")
	iv := ciphertext[:des.BlockSize] // const BlockSize = 8

	mode := cipher.NewCBCEncrypter(block, iv)

	plaintext = PKCS5Padding(plaintext, block.BlockSize())

	encrypted := make([]byte, len(plaintext))
	mode.CryptBlocks(encrypted, plaintext)
	return base64.URLEncoding.EncodeToString(encrypted)

}

func Decrypt(encrypted_ string) string {
	encrypted, _ := base64.URLEncoding.DecodeString(encrypted_)
	//encrypted := []byte(encrypted_)
	block, _ := des.NewTripleDESCipher([]byte(triplekey))
	ciphertext := []byte("abcdef1234567890")
	iv := ciphertext[:des.BlockSize] // const BlockSize = 8

	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(encrypted))
	decrypter.CryptBlocks(decrypted, encrypted)

	decrypted = PKCS5UnPadding(decrypted)

	return string(decrypted)
}

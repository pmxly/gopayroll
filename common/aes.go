package common

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

var key = []byte("k2OJ2i75PBA94TU9")
var iv []byte
/*var iv = []byte("1234567890123456")
func Encrypt(origData []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))

	blockMode.CryptBlocks(crypted, origData)
	return base64.StdEncoding.EncodeToString(crypted), nil
}*/

func Decrypt(crypted string) (string, error) {
	if len(crypted) < 32{
		return "", errors.New("token is invalid")
	}
	iv, _ = hex.DecodeString(crypted[0:32])
	crypted = crypted[32:]
	decodeData,err:=base64.StdEncoding.DecodeString(crypted)
	if err != nil {
		return "",err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(decodeData))
	blockMode.CryptBlocks(origData, decodeData)
	origData = PKCS5UnPadding(origData)
	return string(origData), nil
}

/*func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext) % blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}*/

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length - 1])
	return origData[:(length - unpadding)]
}
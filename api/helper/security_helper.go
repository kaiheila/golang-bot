package helper

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	log "github.com/sirupsen/logrus"
	"strings"
)

func DecryptData(data string, encryptKey string) (error, []byte) {

	rawBase64Decoded, err := base64.StdEncoding.DecodeString(data)

	if err != nil {
		log.Error(err)
		return err, nil
	}
	iv := rawBase64Decoded[:16]
	decryptedContent := string(rawBase64Decoded[16:])
	return Ase256CBCDecode(decryptedContent, []byte(encryptKey), iv)
}

// PKCS7 填充函数
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, len(data)+padding)
	copy(padtext, data)
	for i := len(data); i < len(padtext); i++ {
		padtext[i] = byte(padding)
	}
	return padtext
}

func processPassphrase(passphrase []byte) []byte {
	// 如果密码短语短于 32 字节，用 '\x00' 填充

	if len(passphrase) < 32 {
		encryptKeyUsed := strings.Builder{}
		encryptKeyUsed.Write(passphrase)
		encryptKeyUsed.Write(bytes.Repeat([]byte{byte(0)}, 32-len(passphrase)))
		passphrase = []byte(encryptKeyUsed.String())
	}
	// 如果密码短语长于 32 字节，截断
	if len(passphrase) > 32 {
		passphrase = passphrase[:32]
	}
	return passphrase
}

// 加密函数
func encryptAES256CBC(plaintext []byte, key []byte, iv []byte) ([]byte, error) {
	usedKey := processPassphrase(key)
	block, err := aes.NewCipher(usedKey)
	if err != nil {
		return nil, err
	}

	// 对明文进行填充
	plaintext = pkcs7Padding(plaintext, aes.BlockSize)

	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	return ciphertext, nil
}

// PKCS7 去除填充函数
func pkcs7UnPadding(data []byte) []byte {
	length := len(data)
	padding := int(data[length-1])
	return data[:(length - padding)]
}

// 解密函数
func decryptAES256CBC(ciphertext []byte, key []byte, iv []byte) ([]byte, error) {
	useKey := processPassphrase(key)
	block, err := aes.NewCipher(useKey)
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	// 去除填充
	plaintext = pkcs7UnPadding(plaintext)

	return plaintext, nil
}
func Aes256CBCEncode(plaintext string, encryptKey []byte, iv []byte) (error, []byte) {
	bKey := processPassphrase(encryptKey)
	bPlaintext := PKCS5Padding([]byte(plaintext), aes.BlockSize)
	block, err := aes.NewCipher(bKey)
	if err != nil {
		return err, nil
	}
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, bPlaintext)
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(ciphertext)))
	base64.StdEncoding.Encode(dst, ciphertext)
	return nil, dst
}

func Ase256CBCDecode(cipherText string, encKey []byte, iv []byte) (error, []byte) {

	cipherTextDecoded, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		log.Error(err)
		return err, nil
	}
	usedEncKey := processPassphrase(encKey)
	bKey := usedEncKey
	bIV := iv

	block, err := aes.NewCipher(bKey)
	if err != nil {
		log.Error(err)
		return err, nil
	}

	mode := cipher.NewCBCDecrypter(block, bIV)
	mode.CryptBlocks([]byte(cipherTextDecoded), []byte(cipherTextDecoded))
	return nil, PKCS5Trimming(cipherTextDecoded)
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}

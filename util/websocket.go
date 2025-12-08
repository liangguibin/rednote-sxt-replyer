package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// ECB 模式加密（无IV，分块独立加密）
func ecbEncrypt(dst, src []byte, block cipher.Block) {
	for i := 0; i < len(src); i += block.BlockSize() {
		block.Encrypt(dst[i:], src[i:])
	}
}

// PKCS7 填充
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// AESECBEncrypt AES-ECB 加密
// 将明文和密钥解析为 UTF-8 字节
// 使用 AES-ECB 模式 + PKCS7 填充进行加密
// 输出 Base64 格式的加密结果
func AESECBEncrypt(text string, secret string) (string, error) {
	// 1.将明文和密钥解析为 UTF-8 字节
	plaintext := []byte(text)
	key := []byte(secret)
	// 2.创建 AES 密码实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	// 3.添加 PKCS7 填充
	plaintext = pkcs7Pad(plaintext, aes.BlockSize)
	// 4.创建 ECB 加密器
	ciphertext := make([]byte, len(plaintext))
	ecbEncrypt(ciphertext, plaintext, block)
	// 5.返回 Base64 编码结果
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

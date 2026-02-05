// Package utils 签名工具
package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

// SignData 使用RSA私钥签名数据（SHA1withRSA算法�?func SignData(data string, privateKeyPEM string) (string, error) {
	// 解析私钥
	privateKey, err := ParsePrivateKey(privateKeyPEM)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	// 使用SHA1哈希
	hash := sha1.Sum([]byte(data))

	// 使用RSA PKCS#1 v1.5签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, hash[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign data: %w", err)
	}

	// Base64编码
	return base64.StdEncoding.EncodeToString(signature), nil
}

// SignDataWithRSAKey 使用已解析的RSA私钥签名数据（高性能版本�?func SignDataWithRSAKey(data string, privateKey *rsa.PrivateKey) (string, error) {
	// 使用SHA1哈希
	hash := sha1.Sum([]byte(data))

	// 使用RSA PKCS#1 v1.5签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, hash[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign data: %w", err)
	}

	// Base64编码
	return base64.StdEncoding.EncodeToString(signature), nil
}

// ParsePrivateKey 解析PEM格式的RSA私钥
func ParsePrivateKey(privateKeyPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to decode private key PEM")
	}

	var privateKey *rsa.PrivateKey
	var err error

	// 尝试解析PKCS#1格式
	if privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return privateKey, nil
	}

	// 尝试解析PKCS#8格式
	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		if rsaKey, ok := key.(*rsa.PrivateKey); ok {
			return rsaKey, nil
		}
		return nil, fmt.Errorf("not an RSA private key")
	}

	return nil, fmt.Errorf("failed to parse private key")
}

// VerifySignature 验证RSA签名（用于测试）
func VerifySignature(data, signature, publicKeyPEM string) error {
	// 解析公钥
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return fmt.Errorf("failed to decode public key PEM")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("not an RSA public key")
	}

	// 解码签名
	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %w", err)
	}

	// 计算哈希
	hash := sha1.Sum([]byte(data))

	// 验证签名
	return rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA1, hash[:], signatureBytes)
}

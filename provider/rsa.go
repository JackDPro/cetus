package provider

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"os"
)

type Rsa struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func NewRsa() *Rsa {
	return &Rsa{}
}

func NewRsaByKeyContent(privateKey []byte) (*Rsa, error) {
	r := &Rsa{}
	key, err := r.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return nil, err
	}
	r.PrivateKey = key
	r.PublicKey = &key.PublicKey
	return r, nil
}

func NewRsaByKeyFile(privateKeyPath string) (*Rsa, error) {
	file, err := os.Open(privateKeyPath)
	if err != nil {
		return nil, err
	}
	bytes, err := io.ReadAll(file)
	return NewRsaByKeyContent(bytes)
}

func (r *Rsa) Signature(data []byte) ([]byte, error) {
	hashed := sha256.Sum256(data)
	signature, err := rsa.SignPKCS1v15(rand.Reader, r.PrivateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func (r *Rsa) VerifySignature(data []byte, signature []byte) error {
	hashed := sha256.Sum256(data)
	return rsa.VerifyPKCS1v15(r.PublicKey, crypto.SHA256, hashed[:], signature)
}

func (r *Rsa) ParseRSAPublicKeyFromPEM(keyString []byte) (*rsa.PublicKey, error) {
	// 解码 PEM 格式的公钥
	pemBlock, _ := pem.Decode(keyString)
	if pemBlock == nil {
		return nil, errors.New("error decoding PEM block containing public key")
	}

	// 解析 PEM 编码的公钥为 *rsa.PublicKey
	publicKeyInterface, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing public key %s ", err)
	}

	// 转换为 *rsa.PublicKey
	rsaPublicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("parsed public key is not an RSA public key")
	}
	return rsaPublicKey, nil
}

func (r *Rsa) ParseRSAPrivateKeyFromPEM(keyString []byte) (*rsa.PrivateKey, error) {
	// 解码PEM格式的字符串
	block, _ := pem.Decode(keyString)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}
	// 解析PKCS#8格式的私钥
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 转换为RSA私钥
	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("parsed private key is not an RSA private key")
	}
	return rsaPrivateKey, nil
}

func (r *Rsa) ExportRSAPrivateKeyToPEM(privateKey *rsa.PrivateKey) ([]byte, error) {
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	privateKeyPem := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privateKeyBytes})
	return privateKeyPem, nil
}

func (r *Rsa) ExportRSAPublicKeyToPEM(publicKey *rsa.PublicKey) ([]byte, error) {
	// 将公钥转换为 DER 编码
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	// 创建 PEM 块
	pemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	return pem.EncodeToMemory(pemBlock), nil
}

func (r *Rsa) GeneratePrivateKey(bits int) (*rsa.PrivateKey, error) {
	// 生成 RSA 密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	// 将私钥转换为 PKCS8 编码
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	// 解析 PKCS8 编码的私钥为 *rsa.PrivateKey 类型
	parsedPrivateKey, err := x509.ParsePKCS8PrivateKey(privateKeyBytes)
	if err != nil {
		return nil, err
	}
	rsaPrivateKey, ok := parsedPrivateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("无法转换为 *rsa.PrivateKey 类型")
	}
	return rsaPrivateKey, nil
}

func (r *Rsa) EncodeByPrivateKey(data []byte) ([]byte, error) {
	return nil, nil
}

func (r *Rsa) DecodeByPublicKey(data []byte) ([]byte, error) {
	return nil, nil
}

func (r *Rsa) EncodeByPublicKey(data []byte) ([]byte, error) {
	return nil, nil
}

func (r *Rsa) DecodeByPrivateKey(data []byte) ([]byte, error) {
	return nil, nil
}

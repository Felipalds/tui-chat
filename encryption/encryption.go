package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

var AesKey []byte

func GetPublicKey(publicString string) (*rsa.PublicKey, error) {

	decodedKey, err := base64.StdEncoding.DecodeString(publicString)
	if err != nil {
		return nil, err
	}

	pubKey, err := x509.ParsePKIXPublicKey(decodedKey)
	if err != nil {
		return nil, err
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("chave pública não é do tipo RSA")
	}

	return rsaPubKey, nil
}

func GenerateAESKey() ([]byte, error) {
	key := make([]byte, 32) // AES-256
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func EncryptWithPublicKey(pub *rsa.PublicKey, msg []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, pub, msg)
}

func ReadPublicKey(publicString string) (string, error) {

	rsaPublicKey, err := GetPublicKey(publicString)
	if err != nil {
		return "", err
	}
	AesKey, err = GenerateAESKey()
	if err != nil {
		return "", err
	}
	encryptedAESKey, err := EncryptWithPublicKey(rsaPublicKey, AesKey)
	if err != nil {
		return "", err
	}

	base64Key := base64.StdEncoding.EncodeToString(encryptedAESKey)
	return base64Key, nil

}

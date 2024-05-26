package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
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

	fmt.Println("Chave pública decrypted:", rsaPubKey)

	return rsaPubKey, nil
}

func GenerateAESKey() ([]byte, error) {
	key := make([]byte, 32) // 32 bytes for AES-256
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		panic(err.Error())
	}

	fmt.Println("Random AES key: \n", key)
	return key, nil
}

func EncryptWithPublicKey(pub *rsa.PublicKey, msg []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, pub, msg)
}

func ReadPublicKey(publicString string) (string, []byte, error) {

	publicString = strings.Replace(publicString, "\n", "", -1)
	rsaPublicKey, err := GetPublicKey(publicString)
	if err != nil {
		return "", nil, err
	}
	AesKey, err = GenerateAESKey()
	if err != nil {
		return "", nil, err
	}

	encryptedAESKey, err := EncryptWithPublicKey(rsaPublicKey, AesKey)
	if err != nil {
		return "", nil, err
	}

	fmt.Println(encryptedAESKey)
	base64Key := base64.StdEncoding.EncodeToString(encryptedAESKey)
	return base64Key, AesKey, nil

}

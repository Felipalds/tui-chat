package encryption

import (
	"crypto/aes"
	"crypto/cipher"
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

	publicString = strings.Replace(publicString, "\n", "", -1)
	fmt.Println("Encoded Public Key", publicString)
	fmt.Println("lennn Public Key", len(publicString))
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

	fmt.Println(encryptedAESKey)
	base64Key := base64.StdEncoding.EncodeToString(encryptedAESKey)
	return base64Key, nil

}

// Function to encrypt a message using AES-GCM
func EncryptAES(key []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 12) // GCM standard nonce size is 12 bytes
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nil, nonce, []byte(plaintext), nil)

	// Prepend the nonce to the ciphertext
	ciphertextWithNonce := append(nonce, ciphertext...)

	// Base64 encode the result
	return base64.StdEncoding.EncodeToString(ciphertextWithNonce), nil
}

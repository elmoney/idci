package oaep

import (
	"crypto/rsa"
	"errors"
)

type OaepHandler struct {

	privateKey *rsa.PrivateKey

	publicKey *rsa.PublicKey

	initialized bool
}

// Зашифровать данные
func (h *OaepHandler) Encrypt(data []byte) ([]byte, error) {

	if !h.initialized {
		return nil, errors.New("Handler not initialized!")
	}

	encrypted, err := EncryptBigData(h.publicKey, data)
	if err != nil {
		return nil, err
	}

	return encrypted, nil
}

// Дешифровать данные
func (h *OaepHandler) Decrypt(data []byte) ([]byte, error) {

	if !h.initialized {
		return nil, errors.New("Handler not initialized!")
	}

	decrypted, err := DecryptBigData(h.privateKey, data)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

// Подписать данные
func (h *OaepHandler) Sign(data []byte) ([]byte, error) {

	if !h.initialized {
		return nil, errors.New("Handler not initialized!")
	}

	sign, err := Sign(h.privateKey, data)
	if err != nil {
		return nil, err
	}

	return sign, nil
}

// Верифицировать подпись
func (h *OaepHandler) Verify(data, sign []byte, publicKey *rsa.PublicKey) bool {

	err := Verify(publicKey, data, sign)

	if err != nil {
		return false
	}

	return true
}

// Загрузить публичный ключ набором байт
func (h *OaepHandler) LoadPublicKey(publicKey []byte) error {

	var err error

	pemBlock, err := GetPemBlock(publicKey, false)
	if err != nil {
		h.initialized = false
		return err
	}

	h.publicKey, err = GetPublicKey(pemBlock)
	if err != nil {
		h.publicKey = nil
		h.initialized = false
		return err
	}

	checkInitialization(h)

	return nil
}

// Загрузить публичный ключ из файла
func (h *OaepHandler) LoadPublicKeyFile(path string) error {

	var err error

	pemData, err := GetPemData(path)
	if err != nil {
		h.initialized = false
		return err
	}

	pemBlock, err := GetPemBlock(pemData, false)
	if err != nil {
		h.initialized = false
		return err
	}

	h.publicKey, err = GetPublicKey(pemBlock)
	if err != nil {
		h.publicKey = nil
		h.initialized = false
		return err
	}

	checkInitialization(h)

	return nil
}

func (h *OaepHandler) GetPublicKey() *rsa.PublicKey {

	if h.publicKey != nil {
		return h.publicKey
	}

	return nil
}

// Загрузить приватный ключ набором байт
func (h *OaepHandler) LoadPrivateKey(privateKey []byte) error {

	var err error

	pemBlock, err := GetPemBlock(privateKey, true)
	if err != nil {
		h.initialized = false
		return err
	}

	h.privateKey, err = GetPrivateKey(pemBlock)
	if err != nil {
		h.privateKey = nil
		h.initialized = false
		return err
	}

	checkInitialization(h)

	return nil
}

// Загрузить приватный ключ из файла
func (h *OaepHandler) LoadPrivateKeyFile(path string) error {

	var err error

	pemData, err := GetPemData(path)
	if err != nil {
		h.initialized = false
		return err
	}

	pemBlock, err := GetPemBlock(pemData, true)
	if err != nil {
		h.initialized = false
		return err
	}

	h.privateKey, err = GetPrivateKey(pemBlock)
	if err != nil {
		h.privateKey = nil
		h.initialized = false
		return err
	}

	checkInitialization(h)

	return nil
}

func (h *OaepHandler) Initialized() bool {
	return h.initialized
}

func checkInitialization(h *OaepHandler) {

	if h.publicKey != nil && h.privateKey != nil {
		h.initialized = true
	}
}
package oaep

import (
	"fmt"
	"io/ioutil"
	"crypto/sha256"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"crypto"
	"encoding/binary"
	"errors"
)

type OaepProvider interface {

	LoadPublicKeyFile(path string) error

	LoadPrivateKeyFile(path string) error

	LoadPublicKey(publicKey []byte) error

	GetPublicKey() *rsa.PublicKey

	LoadPrivateKey(privateKey []byte) error

	Encrypt(data []byte) ([]byte, error)

	EncryptBigData(data []byte) ([]byte, error)

	Decrypt(data []byte) ([]byte, error)

	DecryptBigData([]byte) ([]byte, error)

	Sign(data []byte) ([]byte, error)

	Verify(data, sign []byte, publicKey *rsa.PublicKey) bool

	Initialized() bool
}

var (
	oaepLabel = []byte("oaep_provider")
)

// Возвращает pem-формат по переданному ключу id_rsa
func GetPemData(path string) ([]byte, error) {

	var pemData []byte
	var err error

	if pemData, err = ioutil.ReadFile(path); err != nil {
		return nil, err
	}

	return pemData, nil
}

// Возвращает pem-блок переданного pem-формата
func GetPemBlock(pemData []byte, private bool) (*pem.Block, error) {

	var pemBlock *pem.Block
	var blockType = "PUBLIC KEY"
	if private {
		blockType = "RSA PRIVATE KEY"
	}

	if pemBlock, _ = pem.Decode(pemData); pemBlock == nil || pemBlock.Type != blockType {
		fmt.Println(pemBlock.Type)
		return nil, fmt.Errorf("No valid PEM data")
	}

	return pemBlock, nil
}

// Возвращает объект закрытого ключа
func GetPrivateKey(pemBlock *pem.Block) (*rsa.PrivateKey, error) {

	var err error
	var privKey *rsa.PrivateKey

	if privKey, err = x509.ParsePKCS1PrivateKey(pemBlock.Bytes); err != nil {
		return nil, fmt.Errorf("Private key can't be decoded: %s", err)
	}

	return privKey, nil
}

// Возвращает объект закрытого ключа
func GetPublicKey(pemBlock *pem.Block) (*rsa.PublicKey, error) {

	var err error
	var pubKey interface{}

	if pubKey, err = x509.ParsePKIXPublicKey(pemBlock.Bytes); err != nil {
		return nil, fmt.Errorf("Failed to parse RSA public key: %s", err)
	}

	return pubKey.(*rsa.PublicKey), nil
}

func GetPublicKeyFromFile(filePath string) (*rsa.PublicKey, error) {

	pemData, err := GetPemData(filePath)
	if err != nil {
		return nil, err
	}

	pemBlock, err := GetPemBlock(pemData, false)
	if err != nil {
		return nil, err
	}

	return GetPublicKey(pemBlock)
}

func GetPublicKeyFromBytes(publicKey []byte) (*rsa.PublicKey, error) {

	pemBlock, err := GetPemBlock(publicKey, false)
	if err != nil {
		return nil, err
	}

	return GetPublicKey(pemBlock)
}

// Зашифровать публичным ключем данные
func Encrypt(publicKey *rsa.PublicKey, data []byte) ([]byte, error) {

	var err error
	var encrypted []byte

	if encrypted, err = rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, data, oaepLabel); err != nil {
		return nil, err
	}

	return encrypted, nil
}

// Зашифровать большое сообщение(максимально 4294967295 байт)
func EncryptBigData(publicKey *rsa.PublicKey, data []byte) ([]byte, error) {

	var (

		seek = 0
		blockData, seekData, encrypted []byte
		err error
	)

	hash := sha256.New()

	keyLength := (publicKey.N.BitLen() + 7) / 8
	blockLength := keyLength - 2 * hash.Size() - 2
	dataLength := len(data)

	blockData = make([]byte, 4)
	binary.LittleEndian.PutUint32(blockData, uint32(dataLength))

	for {
		if dataLength <= blockLength {
			seekData = data
			seek = dataLength
		} else {


			if seek + blockLength > dataLength {
				seekData = data[seek:]
				seek = dataLength
			} else {
				seekData = data[seek : seek + blockLength]
				seek += blockLength
			}
		}

		if encrypted, err = rsa.EncryptOAEP(hash, rand.Reader, publicKey, seekData, oaepLabel); err != nil {
			return nil, err
		}
		blockData = append(blockData, encrypted...)


		if seek == dataLength {
			break
		}
	}

	return blockData, nil
}

// Дешифровать данные открытым ключом
func Decrypt(privateKey *rsa.PrivateKey, encrypted []byte) ([]byte, error) {

	var err error
	var decrypted []byte

	if decrypted, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encrypted, oaepLabel); err != nil {
		return nil, err
	}

	return decrypted, nil
}

// Дешифровать большое сообщение(максимально 4294967295 байт)
func DecryptBigData(privateKey *rsa.PrivateKey, encrypted []byte) ([]byte, error) {

	var (

		data, decrypted, seekData []byte
		err error
		blocks int
	)

	keyLength := (privateKey.N.BitLen() + 7) / 8

	if len(encrypted) < 4 {
		return nil, errors.New("Incorrect data for decrypt")
	}

	encrypted = encrypted[4:]
	blocks = len(encrypted) / keyLength

	if len(encrypted) % keyLength != 0 {
		return nil, errors.New("Incorrect data for decrypt")
	}

	var i = 0
	for i < blocks {

		data = encrypted[i * keyLength : i * keyLength + keyLength]

		seekData, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, data, oaepLabel)
		if err != nil {
			return nil, err
		}

		decrypted = append(decrypted, seekData...)

		i++
	}

	return decrypted, nil
}

// Подписать данные
func Sign(privateKey *rsa.PrivateKey, data []byte) ([]byte, error) {

	var err error
	var opts rsa.PSSOptions

	newHash := crypto.SHA256

	hash := newHash.New()
	hash.Write(data)
	hashed := hash.Sum(nil)

	opts.SaltLength = rsa.PSSSaltLengthAuto

	sign, err := rsa.SignPSS(rand.Reader, privateKey, newHash, hashed, &opts)
	if err != nil {
		return nil, err
	}

	return sign, err
}

// Верифицировать подпись
func Verify(publicKey *rsa.PublicKey, data, sign []byte) error {

	var opts rsa.PSSOptions

	newHash := crypto.SHA256

	hash := newHash.New()
	hash.Write(data)
	hashed := hash.Sum(nil)

	opts.SaltLength = rsa.PSSSaltLengthAuto

	return rsa.VerifyPSS(publicKey, newHash, hashed, sign, &opts)
}
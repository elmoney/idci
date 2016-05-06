package oaep_test

import (
	"testing"
	"fmt"

	"github.com/idci/oaep"
)

func TestOaepProvider(t *testing.T) {

	var data = []byte("[TEST 1]TEXT DATA")

	// загружаем публичный ключ для шифровки
	pemData, _ := oaep.GetPemData("../openchain/example/chaincode/chaincode_bankid/ssl/webmoney/id_rsa.pub")

	pemBlock, err := oaep.GetPemBlock(pemData, false)
	if err != nil {
		return
	}

	pubKey, _ := oaep.GetPublicKey(pemBlock)

	// загружаем приватный ключ для дешифровки
	pemData, _ = oaep.GetPemData("../openchain/example/chaincode/chaincode_bankid/ssl/webmoney/id_rsa")

	pemBlock, _ = oaep.GetPemBlock(pemData, true)

	privKey, _ := oaep.GetPrivateKey(pemBlock)

	encrypted, err := oaep.Encrypt(pubKey, data)
	if err != nil {
		t.Fatalf("Encrypt error: %s", err)
	}

	fmt.Printf("\nText \"%s\" encrypted to %x\n", data, encrypted)

	decrypted, err := oaep.Decrypt(privKey, encrypted)
	if err != nil {
		t.Fatalf("Decrypt error: %s", err)
	}

	fmt.Printf("\n\nCipher \"%x\" decrypted to text \"%s\"\n", encrypted, decrypted)

	if string(decrypted) != string(data) {
		t.Fatalf("Test not passed: Expecting \"%s\", but was \"%s\"", data, decrypted)
	}
}
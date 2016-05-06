package oaep_test

import (
	"testing"
	"fmt"

	"github.com/idci/oaep"
)

func TestOaepHandler(t *testing.T) {

	var oaepHandler *oaep.OaepHandler
	var data = []byte("[TEST 2]TEXT DATA")

	oaepHandler = &oaep.OaepHandler{}

	oaepHandler.LoadPrivateKeyFile("../openchain/example/chaincode/chaincode_bankid/ssl/webmoney/id_rsa")
	oaepHandler.LoadPublicKeyFile("../openchain/example/chaincode/chaincode_bankid/ssl/webmoney/id_rsa.pub")

	encrypted, err := oaepHandler.Encrypt(data)
	if err != nil {
		t.Fatalf("Encrypt error: %s", err)
	}

	fmt.Printf("\nText \"%s\" encrypted to %x\n", data, encrypted)

	decrypted, err := oaepHandler.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt error: %s", err)
	}

	fmt.Printf("\n\nCipher \"%x\" decrypted to text \"%s\"\n", encrypted, decrypted)

	if string(decrypted) != string(data) {
		t.Fatalf("Test not passed: Expecting \"%s\", but was \"%s\"", data, decrypted)
	}

	sign, err := oaepHandler.Sign(data)
	if err != nil {
		t.Fatalf("Sign error: %s", err)
	}

	fmt.Printf("\nSign from \"%s\": %x\n", data, sign)

	verify := oaepHandler.Verify(data, sign, oaepHandler.GetPublicKey())
	if !verify {
		t.Fatal("Sign not verified!")
		return
	}

	fmt.Println("Sign successfuly verified!")
}
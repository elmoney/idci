package models

import (
	"encoding/json"
	"os"
	"crypto/rsa"
	"github.com/idci/oaep"
)

//Config модель для маппинга настроек
type Config struct {
	Settings     struct {
			     BICertPath     string `json:"bICertPath"`
			     PrivateKeyPath string `json:"privateKeyPath"`
			     PublicKeyPath  string `json:"publicKeyPath"`
			     PeerDeployURL    string `json:"peerDeployUrl"`
			     PeerInvokeURL    string `json:"peerInvokeUrl"`
			     PeerQueryURL     string `json:"peerQueryUrl"`
			     PeerRegistrarURL string `json:"peerRegistrarUrl"`
			     ChainCodeName    string `json:"chainCodeName"`
			     PeerOwnerAlias    string `json:"peerOwnerAlias"`
		             LogPath string `json:"logPath"`
			     TestMode string `json:"testMode"`
		     } `json:"Settings"`
	Globals      struct {
			     ClientPrivateKey *rsa.PrivateKey
			     ClientPublicKey  *rsa.PublicKey
			     BIPublicKey      *rsa.PublicKey
		     }
}

//InitConfig  инициация настроек
func InitConfig() (*Config) {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := Config{}
	decoder.Decode(&configuration)

	pemData, err := oaep.GetPemData(configuration.Settings.PublicKeyPath)

	if (err != nil) {
		panic(err)
	}

	pemBlock, err := oaep.GetPemBlock(pemData, false)

	if (err != nil) {
		panic(err)
	}

	pubKey, err := oaep.GetPublicKey(pemBlock)

	if (err != nil) {
		panic(err)
	}

	pemData, err = oaep.GetPemData(configuration.Settings.PrivateKeyPath)

	if (err != nil) {
		panic(err)
	}

	pemBlock, err = oaep.GetPemBlock(pemData, true)
	if (err != nil) {
		panic(err)
	}

	privKey, err := oaep.GetPrivateKey(pemBlock)
	if (err != nil) {
		panic(err)
	}

	configuration.Globals.ClientPrivateKey = privKey
	configuration.Globals.ClientPublicKey = pubKey

	pemData, err = oaep.GetPemData(configuration.Settings.BICertPath)

	if (err != nil) {
		panic(err)
	}

	pemBlock, err = oaep.GetPemBlock(pemData, false)

	if (err != nil) {
		panic(err)
	}

	pubKey, err = oaep.GetPublicKey(pemBlock)

	if (err != nil) {
		panic(err)
	}

	configuration.Globals.BIPublicKey = pubKey

	return &configuration
}




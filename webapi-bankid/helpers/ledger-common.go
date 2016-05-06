package helpers

import (
	"errors"
	"github.com/idci/webapi-bankid/models"
	"github.com/idci/webapi-bankid/utils"
	"log"
	"net/http"
	"strconv"
)


//InitLedger инициализация леджера
func InitLedger(config *models.Config) (string,error) {

	status, resp := utils.InitCC("init", config.Settings.PeerDeployURL, config.Settings.ChainCodeName)

	if config.Settings.TestMode == "1" {
		utils.InitTestUnits(config.Settings.PeerInvokeURL, config.Settings.ChainCodeName)
	}

	log.Println("init")
	log.Println(status)
	log.Println(resp)

	if http.StatusOK == status {

		return  "",nil

	}

	return  resp,errors.New(strconv.Itoa(status))
}

package helpers

import (
	"github.com/idci/webapi-bankid/utils"
	"log"
	"net/http"
	"errors"
	"strconv"
	"encoding/json"
	"github.com/idci/webapi-bankid/models"
)


//RejectRequestCall отмена идентификации
func RejectRequestCall(config *models.Config,
unitAliasSenderCrypted string,
sign string,
requestID string,
reasonCrypted string) (error) {

	/*
	unitAliasSenderCrypted = args[0]
	sign = args[1]
	requestId = args[2]
	reasonCrypted = args[3]
	*/

	status, resp := utils.InvokeCC(`reject-request`, config.Settings.PeerInvokeURL, config.Settings.ChainCodeName,
		unitAliasSenderCrypted,
		sign,
		requestID,
		reasonCrypted)

	log.Println("reject-request")
	log.Println(status)
	log.Println(resp)

	if (http.StatusOK == status) {
		return nil
	}

	return errors.New(strconv.Itoa(status))

}


//ApproveRequestCall вызов подтверждения Request *** //
func ApproveRequestCall(config *models.Config,
unitAliasSenderCrypted string,
sign string,
requestID string,
urlCrypted string) (error) {

	/*
	unitAliasSenderCrypted = args[0]
	sign = args[1]
	requestId = args[2]
	urlCrypted = args[3]
	*/

	status, resp := utils.InvokeCC(`approve-request`, config.Settings.PeerInvokeURL, config.Settings.ChainCodeName,
		unitAliasSenderCrypted,
		sign,
		requestID,
		urlCrypted)

	log.Println("approve-request")
	log.Println(status)
	log.Println(resp)

	if (http.StatusOK == status) {
		return nil
	}

	return errors.New(strconv.Itoa(status))

}



//CreateRequestCall вызов создания Request *** //
func CreateRequestCall(config *models.Config, unitAliasSenderCrypted string, sign string, requestID string,
clientIDCrypted string, unitAliasRecipientCrypted string, typeSetIdentificationCrypted string, hashPiCrypted string, typeIdentificationCrypted string) (error) {

	/*unitAliasSenderCrypted = args[0]
	sign = args[1]
	requestId = args[2]
	unitAliasRecipientCrypted = args[3]
	typeIdentificationCrypted = args[4]
	clientIdCrypted = args[5]
	typeSetIdentificationCrypted = args[6]
	hashPiCrypted = args[7]*/


	status, resp := utils.InvokeCC(`create-request`, config.Settings.PeerInvokeURL, config.Settings.ChainCodeName,
		unitAliasSenderCrypted,
		sign,
		requestID,
		unitAliasRecipientCrypted,
		typeIdentificationCrypted,
		clientIDCrypted,
		typeSetIdentificationCrypted,
		hashPiCrypted)

	log.Println("create-request")
	log.Println(status)
	log.Println(resp)

	if (http.StatusOK == status) {
		return nil
	}
	return errors.New(strconv.Itoa(status))

}

//Get запрос на идентфикацию
func Get(config *models.Config, ID string, unitAliasSenderCrypted string) (string,error) {

	status, resp := utils.QueryStateWithFunctionGet(config.Settings.ChainCodeName, config.Settings.PeerQueryURL, "get-request", unitAliasSenderCrypted, ID)

	log.Println("get-request")
	log.Println(status)
	log.Println(resp)

	if (http.StatusOK == status) {
		var response models.RequestGetResponse

		err := json.Unmarshal([]byte(resp), &response)

		if (err != nil) {
			log.Println(err)
			return  "",err
		}
		log.Println("request.OK - " + response.OK)

		return response.OK,nil

	}

	return resp, errors.New(strconv.Itoa(status))
}

//RequestsToApproveGet  вытащить все запросы на идентфикацию.
func RequestsToApproveGet(config *models.Config, unitAliasSenderCrypted string) ( string,error) {

	status, resp := utils.QueryStateWithFunctionGet(config.Settings.ChainCodeName, config.Settings.PeerQueryURL, "get-requests-for-approve", unitAliasSenderCrypted)

	log.Println("get-requests-for-approve")
	log.Println(status)
	log.Println(resp)

	if (http.StatusOK == status) {
		var response models.RequestToApproveResponse

		err := json.Unmarshal([]byte(resp), &response)

		if (err != nil) {
			log.Println(err)
			return "",err
		}
		return  response.OK,nil

	}

	return  resp,errors.New(strconv.Itoa(status))

}

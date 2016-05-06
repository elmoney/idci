package controllers

import (
	"log"
	"net/http"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"fmt"
	"crypto/sha256"
	"github.com/idci/webapi-bankid/models"
	"github.com/idci/core/util"
	"github.com/idci/oaep"
	"github.com/idci/webapi-bankid/const/identificationSet"
	"errors"
	"strconv"
	"strings"
	"encoding/hex"
	"encoding/json"
	"github.com/idci/webapi-bankid/helpers"
)


//RequestsToApprove запросы на подтверждление
func RequestsToApprove(r *http.Request, ren render.Render, config *models.Config, params martini.Params)  {
	log.Println(r)
	log.Println(params)

	unitAliasSender, err := oaep.EncryptBigData(config.Globals.BIPublicKey, []byte(config.Settings.PeerOwnerAlias))
	if (err != nil) {
		ren.JSON(http.StatusInternalServerError, err)
		return
	}

	unitAliasSenderCrypted := fmt.Sprintf("%x", unitAliasSender)
	log.Println("unitAliasSenderCrypted - " + unitAliasSenderCrypted)


	desc,err := helpers.RequestsToApproveGet(config, unitAliasSenderCrypted)

	if (err != nil) {
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}

	hexDecoded, err := hex.DecodeString(desc);

	if (err != nil) {
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}

	log.Println(hexDecoded)

	var requests [][]byte

	err = json.Unmarshal(hexDecoded, &requests)

	if (err != nil) {
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}

	var requestPersonValidationList []models.RequestPersonValidation

	for _,element := range requests {
		log.Println(element)
		decrypted, err := oaep.DecryptBigData(config.Globals.ClientPrivateKey, element)
		if (err != nil) {
			ren.JSON(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}
		log.Println(string(decrypted))

		var requestPersonValidation models.RequestPersonValidation
		err = json.Unmarshal(decrypted, &requestPersonValidation)

		if (err != nil) {
			ren.JSON(http.StatusInternalServerError, err)
			log.Println(err)
			return
		}

		requestPersonValidation.Person.HashPersonalInfo =fmt.Sprintf("%x", requestPersonValidation.Person.HashPersonalInfo)

		requestPersonValidationList = append(requestPersonValidationList, requestPersonValidation)
	}
	ren.JSON(http.StatusOK, requestPersonValidationList)
}

//RequestReject  отклонение запроса на идентифкацию
func RequestReject(r *http.Request, ren render.Render, params martini.Params, config *models.Config, request models.RejectRequest) {
	log.Println(r)
	log.Println(request)

	unitAliasSender, err := oaep.EncryptBigData(config.Globals.BIPublicKey, []byte(config.Settings.PeerOwnerAlias))
	if (err != nil) {
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}

	unitAliasSenderCrypted := fmt.Sprintf("%x", unitAliasSender)
	log.Println("unitAliasSenderCrypted - " + unitAliasSenderCrypted)

	signToVerPlain := fmt.Sprintf("%s:%s:reject", config.Settings.PeerOwnerAlias, request.RequestID)

	signedPlainSign, err := oaep.Sign(config.Globals.ClientPrivateKey, []byte(signToVerPlain));
	if (err != nil) {
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}
	sign := fmt.Sprintf("%x", signedPlainSign)
	log.Println("sign - " + sign)

	reason, err := oaep.EncryptBigData(config.Globals.BIPublicKey, []byte(request.Reason))
	if (err != nil) {
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}

	reasonCrypted := fmt.Sprintf("%x", reason)
	log.Println("UrlCrypted - " + reasonCrypted)

	err = helpers.RejectRequestCall(config, unitAliasSenderCrypted, sign, request.RequestID, reasonCrypted)

	if (err != nil) {
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}
	ren.JSON(http.StatusOK, "")
}




//RequestApprove подтверждение запроса на идентифкацию
func RequestApprove(r *http.Request, ren render.Render, params martini.Params, config *models.Config, request models.ApproveRequest) {
	log.Println(r)
	log.Println(request)

	unitAliasSender, err := oaep.EncryptBigData(config.Globals.BIPublicKey, []byte(config.Settings.PeerOwnerAlias))
	if (err != nil) {
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}

	unitAliasSenderCrypted := fmt.Sprintf("%x", unitAliasSender)
	log.Println("unitAliasSenderCrypted - " + unitAliasSenderCrypted)

	signToVerPlain := fmt.Sprintf("%s:%s:approve", config.Settings.PeerOwnerAlias, request.RequestID)

	signedPlainSign, err := oaep.Sign(config.Globals.ClientPrivateKey, []byte(signToVerPlain));
	if (err != nil) {
		ren.JSON(http.StatusInternalServerError, err)
	}
	sign := fmt.Sprintf("%x", signedPlainSign)
	log.Println("sign - " + sign)

	url, err := oaep.EncryptBigData(config.Globals.BIPublicKey, []byte(request.URL))
	if (err != nil) {
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}

	urlCrypted := fmt.Sprintf("%x", url)
	log.Println("UrlCrypted - " + urlCrypted)

	err = helpers.ApproveRequestCall(config, unitAliasSenderCrypted, sign, request.RequestID, urlCrypted)

	if (err != nil) {
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}

	ren.JSON(http.StatusOK, "")

}
//RequestGet получение екрвесте на идентифкацию -- ***
func RequestGet(r *http.Request, ren render.Render, config *models.Config, params martini.Params) {
	log.Println(r)
	log.Println(params)
	requestID := params["id"]

	// зашифрованный отправитель запроса
	unitAliasSender, err := oaep.EncryptBigData(config.Globals.BIPublicKey, []byte(config.Settings.PeerOwnerAlias))
	if (err != nil) {
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}

	unitAliasSenderCrypted := fmt.Sprintf("%x", unitAliasSender)
	log.Println("unitAliasSenderCrypted - " + unitAliasSenderCrypted)

	desc,error := helpers.Get(config, requestID, unitAliasSenderCrypted)

	if (error != nil) {
		ren.JSON(http.StatusInternalServerError, desc)
		log.Println(err)
		return
	}

	hexDecoded, err := hex.DecodeString(desc);

	if (err != nil) {
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}

	decryptedResponse, err := oaep.DecryptBigData(config.Globals.ClientPrivateKey, hexDecoded)

	if (err != nil) {
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}

	var requestPersonValidation models.RequestPersonValidation

	err = json.Unmarshal(decryptedResponse, &requestPersonValidation)

	if (err != nil) {
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}

	requestPersonValidation.Person.HashPersonalInfo =fmt.Sprintf("%x", requestPersonValidation.Person.HashPersonalInfo)
	ren.JSON(http.StatusOK, requestPersonValidation)
}


//RequestCreate Создание запроса от ЗБ
func RequestCreate(r *http.Request, ren render.Render, params martini.Params, config *models.Config, request models.CreateRequest) {
	log.Println(r)
	log.Println(request)

	requestID := util.GenerateUUID()
	log.Println("requestId - " + requestID)

	plainSign := fmt.Sprintf("%s:%s", config.Settings.PeerOwnerAlias, requestID)
	log.Println("plainSign - " + plainSign)


	//подпись запроса
	signedPlainSign, err := oaep.Sign(config.Globals.ClientPrivateKey, []byte(plainSign));
	if (err != nil) {
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}
	sign := fmt.Sprintf("%x", signedPlainSign)
	log.Println("sign - " + sign)


	// зашифрованный отправитель запроса
	unitAliasSender, err := oaep.EncryptBigData(config.Globals.BIPublicKey, []byte(config.Settings.PeerOwnerAlias))
	if (err != nil) {
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}

	unitAliasSenderCrypted := fmt.Sprintf("%x", unitAliasSender)
	log.Println("unitAliasSenderCrypted - " + unitAliasSenderCrypted)


	// зашифрованный идентифкатор клиента в ИБ
	clientID, err := oaep.EncryptBigData(config.Globals.BIPublicKey, []byte(request.ClientID))
	if (err != nil) {
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}
	clientIDCrypted := fmt.Sprintf("%x", clientID)
	log.Println("clientIdCrypted - " + clientIDCrypted)


	// зашифрованный  идентификатор  ИБ
	recipientAlias, err := oaep.EncryptBigData(config.Globals.BIPublicKey, []byte(request.RecipientAlias))
	if (err != nil) {
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}
	unitAliasRecipientCrypted := fmt.Sprintf("%x", recipientAlias)
	log.Println("unitAliasRecipientCrypted - " + unitAliasRecipientCrypted)


	//ApiConsts.BankSet
	//if(request.IdentificationSet.SetType==1)


	hashPlain := ""

	switch request.IdentificationSet.SetType {
	case identificationSet.BankSet:
		{
			hashPlain = fmt.Sprintf("%x",
				strings.ToLower(request.IdentificationSet.BankSet.AccountNumber))
		}
	case identificationSet.Standard:
		{
			hashPlain = fmt.Sprintf("%x;%x;%x;%x;%x",
				strings.ToLower(request.IdentificationSet.PersonalData.LastName),
				strings.ToLower(request.IdentificationSet.PersonalData.FirstName),
				strings.ToLower(request.IdentificationSet.PersonalData.MiddleName),
				strings.ToLower(request.IdentificationSet.PersonalData.PasSer),
				strings.ToLower(request.IdentificationSet.PersonalData.PasNumber))
		}
	default:
		{
			ren.JSON(http.StatusInternalServerError, errors.New("incorrect identificationSet"))

			return
		}
	}

	typeSetIdentification, err := oaep.EncryptBigData(config.Globals.BIPublicKey, []byte(strconv.Itoa(request.IdentificationSet.SetType)))
	if (err != nil) {
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}
	typeSetIdentificationCrypted := fmt.Sprintf("%x", typeSetIdentification)
	log.Println("typeSetIdentificationCrypted - " + typeSetIdentificationCrypted)

	log.Println("hashPlain(hexed) - " + hashPlain)

	hash := sha256.New()
	hash.Write([]byte(hashPlain))
	md := hash.Sum(nil)

	hashPi, err := oaep.EncryptBigData(config.Globals.BIPublicKey, md);
	if (err != nil) {
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}
	hashPiCrypted := fmt.Sprintf("%x", hashPi)
	log.Println("hashPiCrypted - " + hashPiCrypted)


	//typeIdentificationCrypted

	typeIdentification, err := oaep.EncryptBigData(config.Globals.BIPublicKey, []byte(request.TypeIdentification))
	if (err != nil) {
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}
	typeIdentificationCrypted := fmt.Sprintf("%x", typeIdentification)
	log.Println("unitAliasRecipientCrypted - " + typeIdentificationCrypted)


	err = helpers.CreateRequestCall(config, unitAliasSenderCrypted, sign, requestID, clientIDCrypted, unitAliasRecipientCrypted, typeSetIdentificationCrypted, hashPiCrypted,typeIdentificationCrypted)

	if (err != nil) {
		ren.JSON(http.StatusInternalServerError, err)
		log.Println(err)
		return
	}

	ren.JSON(http.StatusOK, requestID)

}
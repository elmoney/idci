package chaincode_call_handlers

import (
	"errors"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"encoding/json"
	"regexp"

	"github.com/idci/core/chaincode/shim"
	"github.com/op/go-logging"
	"github.com/idci/examples/chaincode/go/chaincode_bankid/request"
	"github.com/idci/oaep"
	"github.com/idci/examples/chaincode/go/chaincode_bankid/bankid_constants"
)

type QueryHandler struct {

	log *logging.Logger

	oaepHandler *oaep.OaepHandler
}

func (q *QueryHandler) Init(log *logging.Logger, oaepHandler *oaep.OaepHandler) {

	q.log = log
	q.oaepHandler = oaepHandler
}

func (q *QueryHandler) GetRequestsForApprove(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	if len(args) != 1 {

		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	var (

		senderCrypted string

		senderDecrypted,
		stateBox,
		reqLinkCrypted,
		reqLinkDecrypted,
		requestPersonValidationBytes,
		pubKeySenderBytes,
		hexDecoded []byte

		pubKeySender *rsa.PublicKey

		reqLink	request.RequestLink

		requests [][]byte

		methodName = "GET_REQUESTS_FOR_APPROVE_BY_UNIT"

		err error
	)

	senderCrypted = args[0]

	if hexDecoded, err = hex.DecodeString(senderCrypted); err != nil {

		q.log.Error(bankid_constants.ErrorFormatHexDecodeInternal, methodName, "UNIT_SENDER", err)

		return nil, errors.New(fmt.Sprintf(bankid_constants.ErrorFormatHexDecodeCommon, "UNIT_SENDER"))
	}

	if senderDecrypted, err = q.oaepHandler.Decrypt(hexDecoded); err != nil {

		q.log.Error(bankid_constants.ErrorFormatDecryptObjectInternal, methodName, "UNIT_SENDER", err)

		return nil, errors.New(fmt.Sprintf(bankid_constants.ErrorFormatDecryptObjectCommon, "UNIT_SENDER"))
	}

	// Достаем необходимый контейнер линков
	stateBox, err = stub.GetState(bankid_constants.BoxReqToVerificationUnitName)
	if err != nil {

		q.log.Error(bankid_constants.ErrorFormatGettingState, methodName, err)
		return nil, errors.New("INTERNAL_ERROR")
	}

	// Десериализуем контейнер
	var boxRequest [][]byte
	err = json.Unmarshal(stateBox, &boxRequest)
	if err != nil {
		q.log.Error(bankid_constants.ErrorFormatDeserializeContainer, methodName, err)
		return nil, errors.New("INTERNAL_ERROR")
	}

	if pubKeySenderBytes, err = stub.GetState(
		fmt.Sprintf("%s%s", bankid_constants.UnitPubKeysStatePrefix, string(senderDecrypted))); err != nil {

		q.log.Error(bankid_constants.ErrorFormatGettingState, methodName, err)
		return nil, errors.New("INTERNAL_ERROR")
	}

	if len(pubKeySenderBytes) == 0 {

		q.log.Warning("[%s]Public key not found in state", methodName)

		return nil, errors.New("NOT_FOUND_SENDER_PUBLIC_KEY_ERROR")
	}

	if pubKeySender, err = oaep.GetPublicKeyFromBytes(pubKeySenderBytes); err != nil {

		q.log.Error("[%s]Error at loading public key: %s", methodName, err)
		return nil, errors.New("INTERNAL_ERROR")
	}

	for _, reqLinkCrypted = range boxRequest {

		if reqLinkDecrypted, err = q.oaepHandler.Decrypt(reqLinkCrypted); err != nil {
			q.log.Error(bankid_constants.ErrorFormatDecryptObjectInternal, methodName, err)
			return nil, errors.New("INTERNAL_ERROR")
		}

		if err = json.Unmarshal(reqLinkDecrypted, &reqLink); err != nil {
			q.log.Error(bankid_constants.ErrorFormatDeserializeObject, methodName, err)
			return nil, errors.New("INTERNAL_ERROR")
		}

		if reqLink.Recipient == string(senderDecrypted) {

			if requestPersonValidationBytes, err = stub.GetState(reqLink.Id); err != nil {

				q.log.Error(bankid_constants.ErrorFormatGettingState, methodName, err)

				continue
			}

			// Дешифруем ключом BankId
			requestPersonValidationBytes, err = q.oaepHandler.Decrypt(requestPersonValidationBytes)
			if err != nil {

				q.log.Error(bankid_constants.ErrorFormatDecryptObjectInternal, methodName, err)

				return nil, errors.New("INTERNAL_ERROR")
			}

			// Десериализуем объект
			reqPersonValidation := &request.RequestPersonValidation{}
			err = json.Unmarshal(requestPersonValidationBytes, reqPersonValidation)
			if err != nil {

				q.log.Error(bankid_constants.ErrorFormatDeserializeObject, methodName, err)
				return nil, errors.New("REQUEST_PERSON_VALIDATION_UNMARSHAL_ERROR")
			}

			q.log.Notice("[%s]Request with id \"%s\" - %s\n", methodName, reqLink.Id, reqPersonValidation)

			if requestPersonValidationBytes, err =
			oaep.EncryptBigData(pubKeySender, requestPersonValidationBytes); err != nil {

				q.log.Error(bankid_constants.ErrorFormatEncryptObjectInternal, methodName, err)
			}

			requests = append(requests, requestPersonValidationBytes)
		}
	}

	if len(requests) == 0 {

		return nil, nil
	}

	result, err := json.Marshal(requests)
	if err != nil {

		q.log.Error(bankid_constants.ErrorFormatSerializeObject, methodName, err)
		return nil, errors.New("INTERNAL_ERROR")
	}

	return []byte(hex.EncodeToString(result)), nil
}

func (q *QueryHandler) GetRequestsAll(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	if len(args) != 1 {

		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	var (

		senderCrypted string

		senderDecrypted,
		stateBox,
		reqLinkCrypted,
		reqLinkDecrypted,
		requestPersonValidationBytes,
		pubKeySenderBytes,
		hexDecoded []byte

		pubKeySender *rsa.PublicKey

		reqLink	request.RequestLink

		requests [][]byte

		methodName = "GET_REQUESTS_ALL"

		err error
	)

	senderCrypted = args[0]

	if hexDecoded, err = hex.DecodeString(senderCrypted); err != nil {

		q.log.Error(bankid_constants.ErrorFormatHexDecodeInternal, methodName, "UNIT_SENDER", err)

		return nil, errors.New(fmt.Sprintf(bankid_constants.ErrorFormatHexDecodeCommon, "UNIT_SENDER"))
	}

	if senderDecrypted, err = q.oaepHandler.Decrypt(hexDecoded); err != nil {

		q.log.Error(bankid_constants.ErrorFormatDecryptObjectInternal, methodName, "UNIT_SENDER", err)

		return nil, errors.New(fmt.Sprintf(bankid_constants.ErrorFormatDecryptObjectCommon, "UNIT_SENDER"))
	}

	// Достаем контейнер линков с запросами на идентификацию
	stateBox, err = stub.GetState(bankid_constants.BoxReqToVerificationUnitName)
	if err != nil {

		q.log.Error(bankid_constants.ErrorFormatGettingState, methodName, err)
		return nil, errors.New("INTERNAL_ERROR")
	}

	// Десериализуем контейнер
	var boxRequest [][]byte
	err = json.Unmarshal(stateBox, &boxRequest)
	if err != nil {
		q.log.Error(bankid_constants.ErrorFormatDeserializeContainer, methodName, err)
		return nil, errors.New("INTERNAL_ERROR")
	}

	// Достаем контейнер линков с подтвержденными запросами
	stateBox, err = stub.GetState(bankid_constants.BoxReqApprovedName)
	if err != nil {

		q.log.Error(bankid_constants.ErrorFormatGettingState, methodName, err)
		return nil, errors.New("INTERNAL_ERROR")
	}

	// Десериализуем контейнер
	var boxRequestApproved [][]byte
	err = json.Unmarshal(stateBox, &boxRequestApproved)
	if err != nil {
		q.log.Error(bankid_constants.ErrorFormatDeserializeContainer, methodName, err)
		return nil, errors.New("INTERNAL_ERROR")
	}

	boxRequest = append(boxRequest, boxRequestApproved...)

	// Достаем контейнер линков с отказанными запросами
	stateBox, err = stub.GetState(bankid_constants.BoxReqRejectedName)
	if err != nil {

		q.log.Error(bankid_constants.ErrorFormatGettingState, methodName, err)
		return nil, errors.New("INTERNAL_ERROR")
	}

	// Десериализуем контейнер
	var boxRequestRejected [][]byte
	err = json.Unmarshal(stateBox, &boxRequestRejected)
	if err != nil {
		q.log.Error(bankid_constants.ErrorFormatDeserializeContainer, methodName, err)
		return nil, errors.New("INTERNAL_ERROR")
	}

	boxRequest = append(boxRequest, boxRequestRejected...)

	if pubKeySenderBytes, err = stub.GetState(
		fmt.Sprintf("%s%s", bankid_constants.UnitPubKeysStatePrefix, string(senderDecrypted))); err != nil {

		q.log.Error(bankid_constants.ErrorFormatGettingState, methodName, err)
		return nil, errors.New("INTERNAL_ERROR")
	}

	if pubKeySender, err = oaep.GetPublicKeyFromBytes(pubKeySenderBytes); err != nil {

		q.log.Error("[%s]Error at loading public key: %s", methodName, err)
		return nil, errors.New("INTERNAL_ERROR")
	}

	for _, reqLinkCrypted = range boxRequest {

		if reqLinkDecrypted, err = q.oaepHandler.Decrypt(reqLinkCrypted); err != nil {
			q.log.Error(bankid_constants.ErrorFormatDecryptObjectInternal, methodName, err)
			return nil, errors.New("INTERNAL_ERROR")
		}

		if err = json.Unmarshal(reqLinkDecrypted, &reqLink); err != nil {
			q.log.Error(bankid_constants.ErrorFormatDeserializeObject, methodName, err)
			return nil, errors.New("INTERNAL_ERROR")
		}

		if reqLink.Recipient == string(senderDecrypted) {

			if requestPersonValidationBytes, err = stub.GetState(reqLink.Id); err != nil {

				q.log.Error(bankid_constants.ErrorFormatGettingState, methodName, err)

				continue
			}

			// Дешифруем ключом BankId
			requestPersonValidationBytes, err = q.oaepHandler.Decrypt(requestPersonValidationBytes)
			if err != nil {

				q.log.Error(bankid_constants.ErrorFormatDecryptObjectInternal, methodName, err)

				return nil, errors.New("INTERNAL_ERROR")
			}

			// Десериализуем объект
			reqPersonValidation := &request.RequestPersonValidation{}
			err = json.Unmarshal(requestPersonValidationBytes, reqPersonValidation)
			if err != nil {

				q.log.Error(bankid_constants.ErrorFormatDeserializeObject, methodName, err)
				return nil, errors.New("REQUEST_PERSON_VALIDATION_UNMARSHAL_ERROR")
			}

			q.log.Notice("[%s]Request with id \"%s\" - %s\n", methodName, reqLink.Id, reqPersonValidation)

			if requestPersonValidationBytes, err =
			oaep.EncryptBigData(pubKeySender, requestPersonValidationBytes); err != nil {

				q.log.Error(bankid_constants.ErrorFormatEncryptObjectInternal, methodName, err)
			}

			requests = append(requests, requestPersonValidationBytes)
		}
	}

	if len(requests) == 0 {

		return nil, nil
	}

	result, err := json.Marshal(requests)
	if err != nil {

		q.log.Error(bankid_constants.ErrorFormatSerializeObject, methodName, err)
		return nil, errors.New("INTERNAL_ERROR")
	}

	return []byte(hex.EncodeToString(result)), nil
}

func (q *QueryHandler) GetRequestById(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var (

		unitAliasSenderCrypted,
		requestId string

		hexDecoded,
		unitAliasSenderDecrypted,
		requestBytes []byte

		senderPublicKey *rsa.PublicKey

		reqPersonValidation *request.RequestPersonValidation

		methodName = "QUERY_REQUEST_BY_ID"

		err error
	)

	if len(args) != 2 {

		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	unitAliasSenderCrypted = args[0]
	requestId = args[1]

	match, err := regexp.MatchString("[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89aAbB][a-f0-9]{3}-[a-f0-9]{12}",
		requestId)
	if err != nil || !match {
		return nil, errors.New("Incorrect format of requestId")
	}

	// Декодируем отправителя запроса
	if hexDecoded, err = hex.DecodeString(unitAliasSenderCrypted); err != nil {
		q.log.Error(bankid_constants.ErrorFormatHexDecodeInternal, methodName, "UNIT_ALIAS_SENDER", err)

		return nil, errors.New(fmt.Sprintf(bankid_constants.ErrorFormatHexDecodeCommon, "UNIT_ALIAS_SENDER"))
	}

	// Расшифровываем отправителя запроса
	if unitAliasSenderDecrypted, err = q.oaepHandler.Decrypt(hexDecoded); err != nil {

		q.log.Error(bankid_constants.ErrorFormatDecryptObjectInternal, methodName, err)

		return nil, errors.New("UNIT_ALIAS_SENDER_DECRYPT_ERROR")
	}

	// Достаем публичный ключ отправителя чтобы проверить подпись
	senderPublicKeyBytes, err := stub.GetState(fmt.Sprintf("%s%s", bankid_constants.UnitPubKeysStatePrefix, unitAliasSenderDecrypted))
	if err != nil {

		q.log.Error("[%s]Error at receiving public key from state: %s", methodName, err)

		return nil,  errors.New("GET_SENDER_PUBLIC_KEY_ERROR")
	}

	if len(senderPublicKeyBytes) == 0 {

		q.log.Info("[%s]Public key not found in state", methodName)

		return nil, errors.New("NOT_FOUND_SENDER_PUBLIC_KEY_ERROR")
	}

	if senderPublicKey, err = oaep.GetPublicKeyFromBytes(senderPublicKeyBytes); err != nil {

		q.log.Error("[%s]Error at loading public key: %s", methodName, err)

		return nil, errors.New("LOAD_SENDER_PUBLIC_KEY_ERROR")
	}

	// Ищем запрос в леджере
	if requestBytes, err = stub.GetState(requestId); err != nil {

		q.log.Error(bankid_constants.ErrorFormatGettingState, methodName, err)

		return nil, errors.New("GET_REQUEST_FROM_STATE_ERROR")
	}

	if len(requestBytes) == 0 {

		return nil, nil
	}

	if requestBytes, err = q.oaepHandler.Decrypt(requestBytes); err != nil {

		q.log.Error(bankid_constants.ErrorFormatDecryptObjectInternal, methodName, err)

		return nil, errors.New("REQUEST_PERSON_VALIDATION_DECRYPT_ERROR")
	}

	// Десериализуем объект
	reqPersonValidation = &request.RequestPersonValidation{}
	err = json.Unmarshal(requestBytes, reqPersonValidation)
	if err != nil {

		q.log.Error(bankid_constants.ErrorFormatDeserializeObject, methodName, err)
		return nil, errors.New("REQUEST_PERSON_VALIDATION_UNMARSHAL_ERROR")
	}

	q.log.Notice("[%s]Request with id \"%s\" - %s\n", methodName, requestId, reqPersonValidation)

	// Шифруем ключем отправителя
	if requestBytes, err = oaep.EncryptBigData(senderPublicKey, requestBytes); err != nil {

		q.log.Error(bankid_constants.ErrorFormatEncryptObjectInternal, methodName, err)

		return nil, errors.New("REQUEST_PERSON_VALIDATION_ENCRYPT_ERROR")
	}

	return []byte(hex.EncodeToString(requestBytes)), nil
}

func (q *QueryHandler) GetRequestLink(stub *shim.ChaincodeStub,
	requestId string, stage request.StageResult) (*request.RequestLink, error) {

	var (

		stateBox,
		reqLinkCrypted,
		reqLinkDecrypted []byte

		boxName string

		reqLink *request.RequestLink

		err error
	)

	// Сохраняем имя нужного контейнера с линками на запрос
	switch stage {

	case request.RequestCreate:
		boxName = bankid_constants.BoxReqCreatedName

	case request.RequestToUnitValidation:
		boxName = bankid_constants.BoxReqToVerificationUnitName

	case request.RequestApproved:
		boxName = bankid_constants.BoxReqApprovedName

	case request.RequestRejected:
		boxName = bankid_constants.BoxReqRejectedName
	}

	// Достаем необходимый контейнер линков
	stateBox, err = stub.GetState(boxName)
	if err != nil {
		return nil, err
	}

	// Десериализуем контейнер
	var boxRequest [][]byte
	err = json.Unmarshal(stateBox, &boxRequest)
	if err != nil {
		q.log.Error(bankid_constants.ErrorFormatDeserializeContainer, "GET_REQUEST", err)
		return nil, err
	}

	reqLink = &request.RequestLink{}
	for _, reqLinkCrypted = range boxRequest {

		reqLinkDecrypted, err = q.oaepHandler.Decrypt(reqLinkCrypted)
		if err != nil {
			q.log.Error(bankid_constants.ErrorFormatDecryptObjectInternal, "GET_REQUEST", err)
			return nil, err
		}

		err = json.Unmarshal(reqLinkDecrypted, reqLink)
		if err != nil {
			q.log.Error(bankid_constants.ErrorFormatDeserializeObject, "GET_REQUEST", err)
			return nil, err
		}

		if reqLink.Id == requestId {
			return reqLink, nil
		}
	}

	return nil, nil
}

func (q *QueryHandler) RequestExists(stub *shim.ChaincodeStub, requestId string) bool {

	data, err := stub.GetState(requestId)

	if err != nil {

		panic("[REQUEST_EXISTS]Critical error: error at request to ledger!")
	}

	if len(data) == 0 {

		return false
	} else {

		return true
	}
}
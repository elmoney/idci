/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"os"

	"github.com/idci/core/chaincode/shim"
	"github.com/idci/examples/chaincode/go/chaincode_bankid/bankid_constants"
	"github.com/idci/examples/chaincode/go/chaincode_bankid/chaincode_call_handlers"
	"github.com/idci/examples/chaincode/go/chaincode_bankid/request"
	"github.com/idci/examples/chaincode/go/chaincode_bankid/ssl/genesis"
	"github.com/op/go-logging"
	//"github.com/idci/examples/chaincode/go/chaincode_bankid/ssl/genesis"
	"github.com/idci/examples/chaincode/go/chaincode_bankid/validators"
	"github.com/idci/oaep"
)

var (
	bankIDLog *logging.Logger

	parameterValidaters        []validators.ParameterValidater
	validatorRequestIdentifier *validators.ValidatorRequestIdentifier

	oaepHandler *oaep.OaepHandler

	queryHandler *chaincode_call_handlers.QueryHandler
)

/*
BankIDChainCode реализует бизнес-логику проекта
*/
type BankIDChainCode struct {
}

/*
Init подготавливает chaincode к работе
	- создание списка партнеров
	- создание и добавление в леджер контейнера запросов
*/
func (t *BankIDChainCode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	box, _ := json.Marshal([][]byte{})

	// Инициализируем контейнеры
	stub.PutState(bankid_constants.BoxReqCreatedName, box)
	stub.PutState(bankid_constants.BoxReqToVerificationUnitName, box)
	stub.PutState(bankid_constants.BoxReqApprovedName, box)
	stub.PutState(bankid_constants.BoxReqRejectedName, box)

	// Добавляем публичный ключ BankId в леджер
	stub.PutState(fmt.Sprintf("%s%s", bankid_constants.UnitPubKeysStatePrefix, "bankid"), genesis.PubKeyBankId)

	return nil, nil
}

// Добавление новой идентифицирующей структуры
func (t *BankIDChainCode) unitAdd(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var (
		unitAlias,
		unitPubKey,
		sign string

		unitPubKeyDecoded,
		hexDecoded []byte

		err error

		methodName = "UNIT_ADD"
	)

	if len(args) < 3 {
		return nil, fmt.Errorf("Incorrect number of arguments. Expecting 3a")
	}

	unitAlias = args[0]
	unitPubKey = args[1]
	sign = args[2]

	// Декодируем подпись
	hexDecoded, err = hex.DecodeString(sign)
	if err != nil {

		bankIDLog.Error(fmt.Sprintf(bankid_constants.ErrorFormatHexDecodeInternal, methodName, "SIGN", err))

		return nil, fmt.Errorf(fmt.Sprintf(bankid_constants.ErrorFormatHexDecodeCommon, "SIGN"))
	}

	// Проверяем структуру на существование
	if checkUnitExist(stub, string(unitAlias)) {

		bankIDLog.Error(fmt.Sprintf("[%s]Unit with same alias already exist in ledger", methodName))
		return nil, fmt.Errorf("UNIT_ALIAS_EXIST_ERROR")
	}

	// Декодируем публичный ключ
	unitPubKeyDecoded, err = hex.DecodeString(unitPubKey)
	if err != nil {

		bankIDLog.Error(fmt.Sprintf("[%s]Error at unit public key hex decode: %s", methodName, err))
		return nil, fmt.Errorf("UNIT_PUB_KEY_DECODE_ERROR")
	}

	// Загружаем ключ из файла
	unitPubKeyFromFile, err := oaep.GetPublicKeyFromFile(fmt.Sprintf("./ssl/%s/id_rsa.pub", unitAlias))
	if err != nil {

		bankIDLog.Error(fmt.Sprintf("[%s]Error at loading public key from file: %s", methodName, err))

		return nil, fmt.Errorf("Error at loading public key from file")
	}

	// Верифицируем подпись
	err = oaep.Verify(unitPubKeyFromFile,
		[]byte(fmt.Sprintf("%s:%s", unitAlias, unitPubKeyDecoded)), hexDecoded)
	if err != nil {
		bankIDLog.Error(fmt.Sprintf("sign - %s\ndata - %x\n",
			[]byte(fmt.Sprintf("%s:%s", unitAlias, unitPubKeyDecoded)), hexDecoded))
		bankIDLog.Error(fmt.Sprintf("[%s]Unit sign is not verified: %s", methodName, err))

		return nil, fmt.Errorf("Unit sign is not verified!")
	}

	// Добавляем структуру и ее публичный ключ в леджер
	err = stub.PutState(fmt.Sprintf("%s%s", bankid_constants.UnitPubKeysStatePrefix, unitAlias), unitPubKeyDecoded)
	if err != nil {

		bankIDLog.Error(fmt.Sprintf("[%s]Error at put state: %s", methodName, err))
		return nil, fmt.Errorf("UNIT_ADD_STATE_ERROR")
	}

	bankIDLog.Info("[%s]Public key of [%s] was added", methodName, unitAlias)

	return nil, nil
}

// Инициализация реквеста
// Первый этап когда клиент пришел в банк и передал свои данные банку для последующей его персонификации
func (t *BankIDChainCode) createRequest(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var (
		unitAliasSenderCrypted,
		unitAliasRecipientCrypted,
		sign,
		requestID,
		clientIDCrypted,
		typeIdentificationCrypted,
		typeSetIdentificationCrypted,
		hashPiCrypted string

		unitAliasSenderDecrypted,
		unitAliasRecepientDecrypted,
		clientIDDecrypted,
		typeIdentificationDecrypted,
		typeSetIdentificationDecrypted,
		hashPiDecrypted,
		hexDecoded []byte

		requestPersonValidation *request.RequestPersonValidation
		senderPublicKey         *rsa.PublicKey

		err,
		stageErr error

		methodName = "CREATE_REQUEST"
	)

	if len(args) != 8 {

		bankIDLog.Error(fmt.Sprintf("Incorrect number of arguments(expecting 8) - %v", len(args)))
		return nil, fmt.Errorf("Incorrect number of arguments. Expecting 8")
	}

	unitAliasSenderCrypted = args[0]
	sign = args[1]
	requestID = args[2]
	unitAliasRecipientCrypted = args[3]
	typeIdentificationCrypted = args[4]
	clientIDCrypted = args[5]
	typeSetIdentificationCrypted = args[6]
	hashPiCrypted = args[7]

	// Проверяю формат идентификатора реквеста на валидность
	if err = validatorRequestIdentifier.Validate([]interface{}{requestID, methodName}); err != nil {
		return nil, err
	}

	// Проверяем реквест
	reqExists := queryHandler.RequestExists(stub, requestID)
	if reqExists {

		bankIDLog.Info("[%s]Request with id \"%s\" already exists in state", methodName, requestID)

		return nil, fmt.Errorf("REQUEST_ALREADY_EXISTS")
	}

	requestPersonValidation = &request.RequestPersonValidation{
		Id: requestID,
		Stage: &request.Stage{
			StageResult: request.RequestToUnitValidation,
			Error:       nil,
		},
		Created:      time.Now(),
		LastModified: time.Now(),
	}

	// Декодируем отправителя запроса
	if hexDecoded, err = hex.DecodeString(unitAliasSenderCrypted); err != nil {
		bankIDLog.Info(bankid_constants.ErrorFormatHexDecodeInternal, methodName, "UNIT_ALIAS_SENDER", err)

		stageErr = fmt.Errorf(fmt.Sprintf(bankid_constants.ErrorFormatHexDecodeCommon, "UNIT_ALIAS_SENDER"))
		saveOrUpdateRequestLink(stub, requestPersonValidation.Id, "", "", stageErr, request.RequestCreate)

		return nil, stageErr
	}

	// Расшифровываем отправителя запроса
	if unitAliasSenderDecrypted, err = oaepHandler.Decrypt(hexDecoded); err != nil {

		bankIDLog.Error(fmt.Sprintf(bankid_constants.ErrorFormatDecryptObjectInternal, methodName, err))

		stageErr = fmt.Errorf(fmt.Sprintf("UNIT_ALIAS_SENDER_DECRYPT_ERROR: %s", unitAliasSenderCrypted))
		saveOrUpdateRequestLink(stub, requestPersonValidation.Id, "", "", stageErr, request.RequestCreate)

		return nil, stageErr
	}

	requestPersonValidation.Sender = string(unitAliasSenderDecrypted)

	// Достаем публичный ключ отправителя чтобы проверить подпись
	senderPublicKeyBytes, err := stub.GetState(fmt.Sprintf("%s%s", bankid_constants.UnitPubKeysStatePrefix, unitAliasSenderDecrypted))
	if err != nil {

		bankIDLog.Info("[%s]Error at receiving public key from state: %s", methodName, err)

		stageErr = fmt.Errorf("GET_SENDER_PUBLIC_KEY_ERROR")
		saveOrUpdateRequestLink(stub, requestPersonValidation.Id, requestPersonValidation.Sender,
			"", stageErr, request.RequestCreate)

		return nil, stageErr
	}

	if len(senderPublicKeyBytes) == 0 {

		bankIDLog.Info("[%s]Public key not found in state", methodName)

		stageErr = fmt.Errorf("NOT_FOUND_SENDER_PUBLIC_KEY_ERROR")
		saveOrUpdateRequestLink(stub, requestPersonValidation.Id, requestPersonValidation.Sender,
			"", stageErr, request.RequestCreate)

		return nil, stageErr
	}

	if senderPublicKey, err = oaep.GetPublicKeyFromBytes(senderPublicKeyBytes); err != nil {

		bankIDLog.Info("[%s]Error at loading public key: %s", methodName, err)

		stageErr = fmt.Errorf("LOAD_SENDER_PUBLIC_KEY_ERROR")
		saveOrUpdateRequestLink(stub, requestPersonValidation.Id, requestPersonValidation.Sender,
			"", stageErr, request.RequestCreate)

		return nil, stageErr
	}

	// Декодируем подпись
	if hexDecoded, err = hex.DecodeString(sign); err != nil {
		bankIDLog.Info(bankid_constants.ErrorFormatHexDecodeInternal, methodName, "SIGN", err)

		stageErr = fmt.Errorf(fmt.Sprintf(bankid_constants.ErrorFormatHexDecodeCommon, "SIGN"))
		saveOrUpdateRequestLink(stub, requestPersonValidation.Id, requestPersonValidation.Sender,
			"", stageErr, request.RequestCreate)

		return nil, stageErr
	}

	// Верифицируем подпись отправителя
	signToVer := fmt.Sprintf("%s:%s", unitAliasSenderDecrypted, requestID)
	signVerified := oaepHandler.Verify([]byte(signToVer), hexDecoded, senderPublicKey)
	if !signVerified {
		bankIDLog.Error(fmt.Sprintf("[%s]Sign %s is not verified!", methodName, sign))

		stageErr = fmt.Errorf(fmt.Sprintf("SIGN_IS_NOT_VERIFIED: %s", sign))
		saveOrUpdateRequestLink(stub, requestPersonValidation.Id, requestPersonValidation.Sender,
			"", stageErr, request.RequestCreate)

		return nil, stageErr
	}

	// Декодируем тип идентификации клиента
	hexDecoded, err = hex.DecodeString(typeIdentificationCrypted)
	if err != nil {
		bankIDLog.Info(bankid_constants.ErrorFormatHexDecodeInternal, methodName, "TYPE_IDENTIFICATION", err)

		stageErr = fmt.Errorf(fmt.Sprintf(bankid_constants.ErrorFormatHexDecodeCommon, "TYPE_IDENTIFICATION"))
		saveOrUpdateRequestLink(stub, requestPersonValidation.Id, requestPersonValidation.Sender,
			requestPersonValidation.Recipient, stageErr, request.RequestCreate)

		return nil, stageErr
	}

	// Расшифровываем тип идентификации клиента
	typeIdentificationDecrypted, err = oaepHandler.Decrypt([]byte(hexDecoded))
	if err != nil {
		bankIDLog.Error(fmt.Sprintf("[%s]Error at type identification decrypt: %s", methodName, err))

		stageErr = fmt.Errorf(fmt.Sprintf("TYPE_IDENTIFICATION_DECRYPT_ERROR: %s", hashPiCrypted))
		saveOrUpdateRequestLink(stub, requestPersonValidation.Id, requestPersonValidation.Sender,
			requestPersonValidation.Recipient, stageErr, request.RequestCreate)

		return nil, stageErr
	}

	typeIdentification, err := strconv.Atoi(string(typeIdentificationDecrypted))
	if err != nil {
		bankIDLog.Error(fmt.Sprintf("[%s]Error at convert type identification: %s", methodName, err))

		stageErr = fmt.Errorf(fmt.Sprintf("TYPE_IDENTIFICATION_CONVERT: %s", typeIdentificationDecrypted))
		saveOrUpdateRequestLink(stub, requestPersonValidation.Id, requestPersonValidation.Sender,
			requestPersonValidation.Recipient, stageErr, request.RequestCreate)

		return nil, stageErr
	}

	if typeIdentification > 2 && typeIdentification <= 0 {
		bankIDLog.Error(fmt.Sprintf("[%s]Incorrect value of client's type identification: %s", methodName, err))

		stageErr = fmt.Errorf(fmt.Sprintf("INCORRECT_TYPE_IDENTIFICATION: %v", typeIdentification))
		saveOrUpdateRequestLink(stub, requestPersonValidation.Id, requestPersonValidation.Sender,
			requestPersonValidation.Recipient, stageErr, request.RequestCreate)

		return nil, stageErr
	}

	requestPersonValidation.TypeIdentification = request.TypeIdentification(typeIdentification)

	// Декодируем идентификатор клиента
	hexDecoded, err = hex.DecodeString(clientIDCrypted)
	if err != nil {
		bankIDLog.Info(bankid_constants.ErrorFormatHexDecodeInternal, methodName, "CLIENT_ID", err)

		stageErr = fmt.Errorf(fmt.Sprintf(bankid_constants.ErrorFormatHexDecodeCommon, "CLIENT_ID"))
		saveOrUpdateRequestLink(stub, requestPersonValidation.Id, requestPersonValidation.Sender,
			requestPersonValidation.Recipient, stageErr, request.RequestCreate)

		return nil, stageErr
	}

	// Расшифровываем идентификатор клиента
	clientIDDecrypted, err = oaepHandler.Decrypt([]byte(hexDecoded))
	if err != nil {
		bankIDLog.Error(fmt.Sprintf("[%s]Error at clientId decrypt: %s", "CREATE_REQUEST", err))

		stageErr = fmt.Errorf(fmt.Sprintf("CLIENT_ID_DECRYPT_ERROR: %s", clientIDCrypted))
		saveOrUpdateRequestLink(stub, requestPersonValidation.Id, requestPersonValidation.Sender,
			requestPersonValidation.Recipient, stageErr, request.RequestCreate)

		return nil, stageErr
	}

	// Декодируем получателя запроса
	hexDecoded, err = hex.DecodeString(unitAliasRecipientCrypted)
	if err != nil {
		bankIDLog.Info(bankid_constants.ErrorFormatHexDecodeInternal, methodName, "UNIT_ALIAS_RECIPIENT", err)

		stageErr = fmt.Errorf(fmt.Sprintf(bankid_constants.ErrorFormatHexDecodeCommon, "UNIT_ALIAS_RECIPIENT"))
		saveOrUpdateRequestLink(stub, requestPersonValidation.Id, requestPersonValidation.Sender,
			"", stageErr, request.RequestCreate)

		return nil, stageErr
	}

	// Расшифровываем получателя запроса
	unitAliasRecepientDecrypted, err = oaepHandler.Decrypt([]byte(hexDecoded))
	if err != nil {
		bankIDLog.Error(fmt.Sprintf("[%s]Error at unitAliasRecepient decrypt: %s", methodName, err))

		stageErr = fmt.Errorf(fmt.Sprintf("UNIT_ALIAS_RECEPIENT_DECRYPT_ERROR: %s", unitAliasRecipientCrypted))
		saveOrUpdateRequestLink(stub, requestPersonValidation.Id, requestPersonValidation.Sender,
			"", stageErr, request.RequestCreate)

		return nil, stageErr
	}

	requestPersonValidation.Recipient = string(unitAliasRecepientDecrypted)

	// Декодируем тип набора идентификации клиента
	hexDecoded, err = hex.DecodeString(typeSetIdentificationCrypted)
	if err != nil {
		bankIDLog.Info(bankid_constants.ErrorFormatHexDecodeInternal, methodName, "TYPE_SET_IDENTIFICATION", err)

		stageErr = fmt.Errorf(fmt.Sprintf(bankid_constants.ErrorFormatHexDecodeCommon, "TYPE_SET_IDENTIFICATION"))
		saveOrUpdateRequestLink(stub, requestPersonValidation.Id, requestPersonValidation.Sender,
			requestPersonValidation.Recipient, stageErr, request.RequestCreate)

		return nil, stageErr
	}

	// Расшифровываем тип набора идентификации клиента
	typeSetIdentificationDecrypted, err = oaepHandler.Decrypt([]byte(hexDecoded))
	if err != nil {
		bankIDLog.Error(fmt.Sprintf("[%s]Error at type identification set decrypt: %s", methodName, err))

		stageErr = fmt.Errorf(fmt.Sprintf("TYPE_SET_IDENTIFICATION_DECRYPT_ERROR: %s", hashPiCrypted))
		saveOrUpdateRequestLink(stub, requestPersonValidation.Id, requestPersonValidation.Sender,
			requestPersonValidation.Recipient, stageErr, request.RequestCreate)

		return nil, stageErr
	}

	// Декодируем хеш персональных данных
	hexDecoded, err = hex.DecodeString(hashPiCrypted)
	if err != nil {
		bankIDLog.Info(bankid_constants.ErrorFormatHexDecodeInternal, methodName, "HASH", err)

		stageErr = fmt.Errorf(fmt.Sprintf(bankid_constants.ErrorFormatHexDecodeCommon, "HASH"))
		saveOrUpdateRequestLink(stub, requestPersonValidation.Id, requestPersonValidation.Sender,
			requestPersonValidation.Recipient, stageErr, request.RequestCreate)

		return nil, stageErr
	}

	// Расшифровываем хеш персональных данных
	hashPiDecrypted, err = oaepHandler.Decrypt([]byte(hexDecoded))
	if err != nil {
		bankIDLog.Error(fmt.Sprintf("[%s]Error at hash decrypt: %s", methodName, err))

		stageErr = fmt.Errorf(fmt.Sprintf("HASH_DECRYPT_ERROR: %s", hashPiCrypted))
		saveOrUpdateRequestLink(stub, requestPersonValidation.Id, requestPersonValidation.Sender,
			requestPersonValidation.Recipient, stageErr, request.RequestCreate)

		return nil, stageErr
	}

	typeSetIdentification, err := strconv.Atoi(string(typeIdentificationDecrypted))
	if err != nil {
		bankIDLog.Error(fmt.Sprintf("[%s]Error at convert type set identification: %s", methodName, err))

		stageErr = fmt.Errorf(fmt.Sprintf("TYPE_SET_DENTIFICATION_CONVERT: %s", typeSetIdentificationDecrypted))
		saveOrUpdateRequestLink(stub, requestPersonValidation.Id, requestPersonValidation.Sender,
			requestPersonValidation.Recipient, stageErr, request.RequestCreate)

		return nil, stageErr
	}

	if typeSetIdentification > 2 && typeSetIdentification <= 0 {
		bankIDLog.Error(fmt.Sprintf("[%s]Incorrect value of client's type set identification: %s", methodName, err))

		stageErr = fmt.Errorf(fmt.Sprintf("INCORRECT_TYPE_SET_IDENTIFICATION: %v", typeSetIdentification))
		saveOrUpdateRequestLink(stub, requestPersonValidation.Id, requestPersonValidation.Sender,
			requestPersonValidation.Recipient, stageErr, request.RequestCreate)

		return nil, stageErr
	}

	requestPersonValidation.Person = &request.Person{
		ClientId:              string(clientIDDecrypted),
		TypeSetIdentification: request.TypeSetIdentification(typeSetIdentification),
		HashPersonalInfo:      string(hashPiDecrypted),
	}

	requestPersonValidationBytes, err := json.Marshal(requestPersonValidation)
	if err != nil {
		bankIDLog.Error(fmt.Sprintf(bankid_constants.ErrorFormatSerializeObject, methodName, err))

		stageErr = fmt.Errorf("REQUEST_PERSON_OBJECT_SERIALIZATION_ERROR")
		saveOrUpdateRequestLink(stub, requestPersonValidation.Id, requestPersonValidation.Sender,
			requestPersonValidation.Recipient, stageErr, request.RequestCreate)

		return nil, stageErr
	}

	// Шифруем все ключом BankId и сохраняем в state
	requestPersonValidationBytes, err = oaepHandler.Encrypt(requestPersonValidationBytes)
	if err != nil {

		bankIDLog.Error(fmt.Sprintf("[%s]Error at encrypt request: %s",
			methodName, err))

		stageErr = fmt.Errorf("REQUEST_PERSON_VALIDATION_ENCRYPT_ERROR")
		requestPersonValidation.Stage.Error = stageErr
		saveOrUpdateRequestLink(stub, requestPersonValidation.Id, requestPersonValidation.Sender,
			requestPersonValidation.Recipient, stageErr, request.RequestCreate)

		return nil, stageErr
	}

	// Сохраняем линк
	err = saveOrUpdateRequestLink(stub, requestPersonValidation.Id, requestPersonValidation.Sender,
		requestPersonValidation.Recipient, nil, request.RequestToUnitValidation)
	if err != nil {
		bankIDLog.Error(fmt.Sprintf("[%s]Error at saving container requests: %s", methodName, err))

		return nil, nil
	}

	// Сохраняем сам реквест
	err = stub.PutState(requestPersonValidation.Id, requestPersonValidationBytes)
	if err != nil {
		bankIDLog.Error(fmt.Sprintf(bankid_constants.ErrorFormatSavingState, methodName, err))
		return nil, nil
	}

	bankIDLog.Info("[%s]Request with id \"%s\" was added to state", methodName, requestPersonValidation.Id)

	return nil, nil
}

func (t *BankIDChainCode) approveRequest(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var (
		requestID,
		sign,
		unitAliasSenderCrypted,
		urlCrypted string

		unitAliasSenderDecrypted,
		urlDecrypted,
		requestPersonValidationBytes,
		hexDecoded []byte

		requestPersonValidation *request.RequestPersonValidation

		methodName = "APPROVE_REQUEST"

		err,
		stageErr error
	)

	if len(args) != 4 {

		bankIDLog.Error(fmt.Sprintf("Incorrect number of arguments(expecting 4) - %v", len(args)))
		return nil, fmt.Errorf("Incorrect number of arguments. Expecting 4")
	}

	unitAliasSenderCrypted = args[0]
	sign = args[1]
	requestID = args[2]
	urlCrypted = args[3]

	// Проверяю формат идентификатора реквеста на валидность
	if err = validatorRequestIdentifier.Validate([]interface{}{requestID, methodName}); err != nil {
		return nil, err
	}

	// Проверяю реквест на существование и правильный статус
	reqLink, err := getRequestLink(stub, requestID, request.RequestToUnitValidation)
	if err != nil {
		bankIDLog.Info("[%s]Error at receiving request from state: %s", methodName, err)

		return nil, fmt.Errorf(fmt.Sprintf("[%s]Error at receiving request \"%s\" from state", methodName,
			requestID))
	}

	if reqLink == nil {
		bankIDLog.Info("[%s]Request with id \"%s\" does not exist in state or has other status",
			methodName, requestID)

		return nil, fmt.Errorf(fmt.Sprintf(`[%s]Request with id \"%s\"
			does not exist in state or has other status`, methodName, requestID))
	}

	// Декодируем отправителя запроса
	if hexDecoded, err = hex.DecodeString(unitAliasSenderCrypted); err != nil {
		bankIDLog.Info(bankid_constants.ErrorFormatHexDecodeInternal, methodName, "UNIT_ALIAS_SENDER", err)

		stageErr = fmt.Errorf(fmt.Sprintf(bankid_constants.ErrorFormatHexDecodeCommon, "UNIT_ALIAS_SENDER"))
		saveOrUpdateRequestLink(stub, reqLink.Id, reqLink.Sender, reqLink.Recipient, stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	// Расшифровываем отправителя запроса
	if unitAliasSenderDecrypted, err = oaepHandler.Decrypt(hexDecoded); err != nil {

		bankIDLog.Error(fmt.Sprintf("[%s]Error at unitAliasSender decrypt: %s", methodName, err))

		stageErr = fmt.Errorf(fmt.Sprintf("UNIT_ALIAS_SENDER_DECRYPT_ERROR: %s", unitAliasSenderCrypted))
		saveOrUpdateRequestLink(stub, reqLink.Id, reqLink.Sender, reqLink.Recipient, stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	// Достаем публичный ключ отправителя чтобы проверить подпись
	senderPublicKeyBytes, err := stub.GetState(fmt.Sprintf("%s%s", bankid_constants.UnitPubKeysStatePrefix, unitAliasSenderDecrypted))
	if err != nil {

		bankIDLog.Info("[%s]Error at receiving public key: %s", methodName, err)

		stageErr = fmt.Errorf("GET_SENDER_PUBLIC_KEY_ERROR")
		saveOrUpdateRequestLink(stub, reqLink.Id, reqLink.Sender,
			reqLink.Recipient, stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	if len(senderPublicKeyBytes) == 0 {

		bankIDLog.Info("[%s]Public key not found in state", methodName)

		stageErr = fmt.Errorf("NOT_FOUND_SENDER_PUBLIC_KEY_ERROR")
		saveOrUpdateRequestLink(stub, reqLink.Id, reqLink.Sender,
			reqLink.Recipient, stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	senderPublicKey, err := oaep.GetPublicKeyFromBytes(senderPublicKeyBytes)
	if err != nil {

		bankIDLog.Info("[%s]Error at loading public key from state: %s", methodName, err)

		stageErr = fmt.Errorf("LOAD_SENDER_PUBLIC_KEY_ERROR")
		saveOrUpdateRequestLink(stub, reqLink.Id, reqLink.Sender,
			"", stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	// Декодируем подпись
	if hexDecoded, err = hex.DecodeString(sign); err != nil {
		bankIDLog.Info(bankid_constants.ErrorFormatHexDecodeInternal, methodName, "SIGN", err)

		stageErr = fmt.Errorf(fmt.Sprintf(bankid_constants.ErrorFormatHexDecodeCommon, "SIGN"))
		saveOrUpdateRequestLink(stub, reqLink.Id, reqLink.Sender,
			reqLink.Recipient, stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	// Верифицируем подпись отправителя
	signToVer := fmt.Sprintf("%s:%s:approve", unitAliasSenderDecrypted, requestID)
	signVerified := oaepHandler.Verify([]byte(signToVer), hexDecoded, senderPublicKey)
	if !signVerified {
		bankIDLog.Error(fmt.Sprintf("[%s]Sign %s is not verified!", methodName, sign))

		stageErr = fmt.Errorf(fmt.Sprintf("SIGN_NOT_VERIFIED: %s", sign))
		saveOrUpdateRequestLink(stub, reqLink.Id, reqLink.Sender,
			"", stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	// Декодируем url со сканом
	if hexDecoded, err = hex.DecodeString(urlCrypted); err != nil {
		bankIDLog.Info(bankid_constants.ErrorFormatHexDecodeInternal, methodName, "URL", err)

		stageErr = fmt.Errorf(fmt.Sprintf(bankid_constants.ErrorFormatHexDecodeCommon, "URL"))
		saveOrUpdateRequestLink(stub, reqLink.Id, reqLink.Sender, reqLink.Recipient, stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	// Расшифровываем url со сканом
	if urlDecrypted, err = oaepHandler.Decrypt(hexDecoded); err != nil {

		bankIDLog.Error(fmt.Sprintf("[%s]Error at url decrypt: %s", methodName, err))

		stageErr = fmt.Errorf(fmt.Sprintf("URL_DECRYPT_ERROR: %s", urlCrypted))
		saveOrUpdateRequestLink(stub, reqLink.Id, reqLink.Sender, reqLink.Recipient, stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	// Проверяем нужная ли структура отвечает
	if reqLink.Recipient != string(unitAliasSenderDecrypted) {

		bankIDLog.Error(fmt.Sprintf("[%s]Wrong unit approve - %s. Expecting - %s", methodName, string(unitAliasSenderDecrypted),
			reqLink.Recipient))

		stageErr = fmt.Errorf(fmt.Sprintf("UNSUITABLE_REQUEST_FOR: %s", urlCrypted))
		saveOrUpdateRequestLink(stub, reqLink.Id, reqLink.Sender, reqLink.Recipient, stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	if requestPersonValidationBytes, err = stub.GetState(reqLink.Id); err != nil {

		bankIDLog.Error(fmt.Sprintf(bankid_constants.ErrorFormatGettingState, methodName, err))

		stageErr = fmt.Errorf("GET_OBJECT_FROM_LEDGER_ERROR")
		saveOrUpdateRequestLink(stub, reqLink.Id, reqLink.Sender, reqLink.Recipient, stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	if requestPersonValidationBytes, err = oaepHandler.Decrypt(requestPersonValidationBytes); err != nil {

		bankIDLog.Error(fmt.Sprintf(bankid_constants.ErrorFormatDecryptObjectInternal, methodName, err))

		stageErr = fmt.Errorf("REQUEST_PERSON_VALIDATION_DECRYPT_ERROR")
		saveOrUpdateRequestLink(stub, reqLink.Id, reqLink.Sender, reqLink.Recipient, stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	requestPersonValidation = &request.RequestPersonValidation{}
	if err = json.Unmarshal(requestPersonValidationBytes, requestPersonValidation); err != nil {

		bankIDLog.Error(fmt.Sprintf(bankid_constants.ErrorFormatDeserializeObject, methodName, err))

		stageErr = fmt.Errorf("REQUEST_PERSON_VALIDATION_UNMARSHAL_ERROR")
		saveOrUpdateRequestLink(stub, reqLink.Id, reqLink.Sender, reqLink.Recipient,
			stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	if requestPersonValidation.Id == "" {

		bankIDLog.Error(fmt.Sprintf("[%s]Request is nil", methodName))

		stageErr = fmt.Errorf("REQUEST_PERSON_VALIDATION_NIL_ERROR")
		saveOrUpdateRequestLink(stub, reqLink.Id, reqLink.Sender, reqLink.Recipient,
			stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	requestPersonValidation.LastModified = time.Now()
	requestPersonValidation.Stage = &request.Stage{
		StageResult: request.RequestApproved,
		Error:       nil,
	}
	requestPersonValidation.SetResult(&request.ApprovedData{
		Url: string(urlDecrypted),
	})

	requestPersonValidationBytes, err = json.Marshal(requestPersonValidation)
	if err != nil {
		bankIDLog.Error(fmt.Sprintf(bankid_constants.ErrorFormatSerializeObject, methodName, err))

		stageErr = fmt.Errorf("REQUEST_PERSON_VALIDATION_OBJECT_SERIALIZATION_ERROR")
		saveOrUpdateRequestLink(stub, requestPersonValidation.Id, requestPersonValidation.Sender,
			requestPersonValidation.Recipient, stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	// Шифруем все ключом BankId и сохраняем в state
	requestPersonValidationBytes, err = oaepHandler.Encrypt(requestPersonValidationBytes)
	if err != nil {

		bankIDLog.Error(fmt.Sprintf("[%s]Error at encrypt request: %s",
			methodName, err))

		stageErr = fmt.Errorf("REQUEST_PERSON_VALIDATION_ENCRYPT_ERROR")
		requestPersonValidation.Stage.Error = stageErr
		saveOrUpdateRequestLink(stub, requestPersonValidation.Id, requestPersonValidation.Sender,
			requestPersonValidation.Recipient, stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	// Переносим линк на реквест в другой контейнер
	err = stageRequestLink(stub, reqLink.Id, request.RequestToUnitValidation, request.RequestApproved)
	if err != nil {
		bankIDLog.Error(fmt.Sprintf(bankid_constants.ErrorFormatSavingState, methodName, err))

		return nil, nil
	}

	// Сохраняем сам реквест
	err = stub.PutState(reqLink.Id, requestPersonValidationBytes)
	if err != nil {
		bankIDLog.Error(fmt.Sprintf(bankid_constants.ErrorFormatSavingState, methodName, err))
		return nil, nil
	}

	bankIDLog.Info("[%s]Request with id \"%s\" was approved", methodName, requestID)

	return nil, nil
}

func (t *BankIDChainCode) rejectRequest(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var (
		requestID,
		sign,
		unitAliasSenderCrypted,
		reasonCrypted string

		unitAliasSenderDecrypted,
		reasonDecrypted,
		requestPersonValidationBytes,
		hexDecoded []byte

		requestPersonValidation *request.RequestPersonValidation

		methodName = "REJECT_REQUEST"

		err,
		stageErr error
	)

	if len(args) != 4 {

		bankIDLog.Error(fmt.Sprintf("Incorrect number of arguments(expecting 4) - %v", len(args)))
		return nil, fmt.Errorf("Incorrect number of arguments. Expecting 4")
	}

	unitAliasSenderCrypted = args[0]
	sign = args[1]
	requestID = args[2]
	reasonCrypted = args[3]

	// Проверяю формат идентификатора реквеста на валидность
	if err = validatorRequestIdentifier.Validate([]interface{}{requestID, methodName}); err != nil {
		return nil, err
	}

	// Проверяю реквест на существование и правильный статус
	reqLink, err := getRequestLink(stub, requestID, request.RequestToUnitValidation)
	if err != nil {
		bankIDLog.Info("[%s]Error at receiving request from state: %s", methodName, err)

		return nil, fmt.Errorf(fmt.Sprintf("[%s]Error at receiving request \"%s\" from state", methodName,
			requestID))
	}

	if reqLink == nil {
		bankIDLog.Info("[%s]Request with id \"%s\" does not exist in state or has other status",
			methodName, requestID)

		return nil, fmt.Errorf(fmt.Sprintf(`[%s]Request with id \"%s\"
			does not exist in state or has other status`, methodName,
			requestID))
	}

	// Декодируем отправителя запроса
	if hexDecoded, err = hex.DecodeString(unitAliasSenderCrypted); err != nil {
		bankIDLog.Info(bankid_constants.ErrorFormatHexDecodeInternal, methodName, "UNIT_ALIAS_SENDER", err)

		stageErr = fmt.Errorf(fmt.Sprintf(bankid_constants.ErrorFormatHexDecodeCommon, "UNIT_ALIAS_SENDER"))
		saveOrUpdateRequestLink(stub, reqLink.Id, reqLink.Sender, reqLink.Recipient, stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	// Расшифровываем отправителя запроса
	if unitAliasSenderDecrypted, err = oaepHandler.Decrypt(hexDecoded); err != nil {

		bankIDLog.Error(fmt.Sprintf("[%s]Error at unitAliasSender decrypt: %s", methodName, err))

		stageErr = fmt.Errorf(fmt.Sprintf("UNIT_ALIAS_SENDER_DECRYPT_ERROR: %s", unitAliasSenderCrypted))
		saveOrUpdateRequestLink(stub, reqLink.Id, reqLink.Sender, reqLink.Recipient, stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	// Достаем публичный ключ отправителя чтобы проверить подпись
	senderPublicKeyBytes, err := stub.GetState(fmt.Sprintf("%s%s", bankid_constants.UnitPubKeysStatePrefix, unitAliasSenderDecrypted))
	if err != nil {

		bankIDLog.Info("[%s]Error at receiving public key: %s", methodName, err)

		stageErr = fmt.Errorf("GET_SENDER_PUBLIC_KEY_ERROR")
		saveOrUpdateRequestLink(stub, reqLink.Id, reqLink.Sender,
			reqLink.Recipient, stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	if len(senderPublicKeyBytes) == 0 {

		bankIDLog.Info("[%s]Public key not found in state", methodName)

		stageErr = fmt.Errorf("NOT_FOUND_SENDER_PUBLIC_KEY_ERROR")
		saveOrUpdateRequestLink(stub, reqLink.Id, reqLink.Sender,
			reqLink.Recipient, stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	senderPublicKey, err := oaep.GetPublicKeyFromBytes(senderPublicKeyBytes)
	if err != nil {

		bankIDLog.Info("[%s]Error at loading public key from state: %s", methodName, err)

		stageErr = fmt.Errorf("LOAD_SENDER_PUBLIC_KEY_ERROR")
		saveOrUpdateRequestLink(stub, reqLink.Id, reqLink.Sender,
			"", stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	// Декодируем подпись
	if hexDecoded, err = hex.DecodeString(sign); err != nil {
		bankIDLog.Info(bankid_constants.ErrorFormatHexDecodeInternal, methodName, "SIGN", err)

		stageErr = fmt.Errorf(fmt.Sprintf(bankid_constants.ErrorFormatHexDecodeCommon, "SIGN"))
		saveOrUpdateRequestLink(stub, reqLink.Id, reqLink.Sender,
			reqLink.Recipient, stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	// Верифицируем подпись отправителя
	signToVer := fmt.Sprintf("%s:%s:reject", unitAliasSenderDecrypted, requestID)
	signVerified := oaepHandler.Verify([]byte(signToVer), hexDecoded, senderPublicKey)
	if !signVerified {
		bankIDLog.Error(fmt.Sprintf("[%s]Sign %s not verified!", methodName, sign))

		stageErr = fmt.Errorf(fmt.Sprintf("SIGN_NOT_VERIFIED: %s", sign))
		saveOrUpdateRequestLink(stub, reqLink.Id, reqLink.Sender,
			"", stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	// Декодируем причину отказа
	if hexDecoded, err = hex.DecodeString(reasonCrypted); err != nil {
		bankIDLog.Info(bankid_constants.ErrorFormatHexDecodeInternal, methodName, "REASON", err)

		stageErr = fmt.Errorf(fmt.Sprintf(bankid_constants.ErrorFormatHexDecodeCommon, "REASON"))
		saveOrUpdateRequestLink(stub, reqLink.Id, reqLink.Sender, reqLink.Recipient, stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	// Расшифровываем причину отказа
	if reasonDecrypted, err = oaepHandler.Decrypt(hexDecoded); err != nil {

		bankIDLog.Error(fmt.Sprintf(bankid_constants.ErrorFormatDecryptObjectInternal, methodName, err))

		stageErr = fmt.Errorf(fmt.Sprintf("URL_DECRYPT_ERROR: %s", reasonCrypted))
		saveOrUpdateRequestLink(stub, reqLink.Id, reqLink.Sender, reqLink.Recipient, stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	// Проверяем нужная ли структура отвечает
	if reqLink.Recipient != string(unitAliasSenderDecrypted) {

		stageErr = fmt.Errorf(fmt.Sprintf("UNSUITABLE_REQUEST_FOR: %s", unitAliasSenderCrypted))
		saveOrUpdateRequestLink(stub, reqLink.Id, reqLink.Sender, reqLink.Recipient, stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	if requestPersonValidationBytes, err = stub.GetState(reqLink.Id); err != nil {

		stageErr = fmt.Errorf("GET_OBJECT_FROM_LEDGER_ERROR")
		saveOrUpdateRequestLink(stub, reqLink.Id, reqLink.Sender, reqLink.Recipient, stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	if requestPersonValidationBytes, err = oaepHandler.Decrypt(requestPersonValidationBytes); err != nil {

		stageErr = fmt.Errorf("REQUEST_PERSON_VALIDATION_DECRYPT_ERROR")
		saveOrUpdateRequestLink(stub, reqLink.Id, reqLink.Sender, reqLink.Recipient, stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	requestPersonValidation = &request.RequestPersonValidation{}
	if err = json.Unmarshal(requestPersonValidationBytes, requestPersonValidation); err != nil {

		bankIDLog.Error(fmt.Sprintf(bankid_constants.ErrorFormatDeserializeObject, methodName, err))

		stageErr = fmt.Errorf("REQUEST_PERSON_VALIDATION_UNMARSHAL_ERROR")
		saveOrUpdateRequestLink(stub, reqLink.Id, reqLink.Sender, reqLink.Recipient,
			stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	if requestPersonValidation.Id == "" {

		bankIDLog.Error(fmt.Sprintf("[%s]Request is nil", methodName))

		stageErr = fmt.Errorf("REQUEST_PERSON_VALIDATION_NIL_ERROR")
		saveOrUpdateRequestLink(stub, reqLink.Id, reqLink.Sender, reqLink.Recipient,
			stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	requestPersonValidation.LastModified = time.Now()
	requestPersonValidation.Stage = &request.Stage{
		StageResult: request.RequestRejected,
		Error:       nil,
	}
	requestPersonValidation.SetResult(&request.RejectedData{
		Reason: string(reasonDecrypted),
	})

	requestPersonValidationBytes, err = json.Marshal(requestPersonValidation)
	if err != nil {
		bankIDLog.Error(fmt.Sprintf(bankid_constants.ErrorFormatSerializeObject, methodName, err))

		stageErr = fmt.Errorf("REQUEST_PERSON_VALIDATION_OBJECT_SERIALIZATION_ERROR")
		saveOrUpdateRequestLink(stub, requestPersonValidation.Id, requestPersonValidation.Sender,
			requestPersonValidation.Recipient, stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	// Шифруем все ключом BankId и сохраняем в state
	requestPersonValidationBytes, err = oaepHandler.Encrypt(requestPersonValidationBytes)
	if err != nil {

		bankIDLog.Error(fmt.Sprintf("[%s]Error at encrypt request: %s",
			methodName, err))

		stageErr = fmt.Errorf("REQUEST_PERSON_VALIDATION_ENCRYPT_ERROR")
		requestPersonValidation.Stage.Error = stageErr
		saveOrUpdateRequestLink(stub, requestPersonValidation.Id, requestPersonValidation.Sender,
			requestPersonValidation.Recipient, stageErr, request.RequestToUnitValidation)

		return nil, stageErr
	}

	// Переносим линк на реквест в другой контейнер
	err = stageRequestLink(stub, reqLink.Id, request.RequestToUnitValidation, request.RequestRejected)
	if err != nil {
		bankIDLog.Error(fmt.Sprintf(bankid_constants.ErrorFormatSavingState, methodName, err))

		return nil, nil
	}

	// Сохраняем сам реквест
	err = stub.PutState(reqLink.Id, requestPersonValidationBytes)
	if err != nil {
		bankIDLog.Error(fmt.Sprintf(bankid_constants.ErrorFormatSavingState, methodName, err))
		return nil, nil
	}

	bankIDLog.Info("[%s]Request with id \"%s\" was rejected", methodName, requestID)

	return nil, nil
}

// Deletes an entity from state
func (t *BankIDChainCode) delete(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("Incorrect number of arguments. Expecting 3")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return nil, fmt.Errorf("Failed to delete state")
	}

	return nil, nil
}

// Invoke занимается вызовом методов chaincode
func (t *BankIDChainCode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	// Handle different functions
	if function == "init" {
		// Инициализация леджера
		return t.Init(stub, function, args)
	} else if function == "unit-add" {
		// Добавление новой идентифицируещей структуры
		return t.unitAdd(stub, args)
	} else if function == "create-request" {
		// Создание реквеста с клиентским идентификатором, хешем и идентификатором структуры его
		// 	персонифицировавшим
		return t.createRequest(stub, args)
	} else if function == "approve-request" {
		// Успешная проверка идентифицирующей структурой и отправка запроса BankId
		return t.approveRequest(stub, args)
	} else if function == "reject-request" {
		// Отказ в идентификации клиента идентифицирующей структурой
		// 	- либо хеш не совпал, с существующими клиентами
		//	- либо у идентифицирующей структуры отстутствует клиент с таким ClientId
		//	- либо клиент не дал согласие на идентификацию
		return t.rejectRequest(stub, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	}

	return nil, fmt.Errorf("Received unknown function invocation")
}

// Query callback representing the query of a chaincode
func (t *BankIDChainCode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	switch function {

	case "get-pubkey":

		unitAlias := args[0]
		if unitAlias == "" {
			return nil, fmt.Errorf("Invalid unit name.")
		}

		pubKey, _ := stub.GetState(fmt.Sprintf("%s%s", bankid_constants.UnitPubKeysStatePrefix, unitAlias))

		return pubKey, nil

	case "get-request":

		return queryHandler.GetRequestById(stub, args)

	case "get-requests-for-approve":

		return queryHandler.GetRequestsForApprove(stub, args)

	case "get-requests-all":

		return queryHandler.GetRequestsAll(stub, args)

	default:
		return nil, fmt.Errorf("Invalid query function name.")
	}
}

func saveOrUpdateRequestLink(stub *shim.ChaincodeStub,
	id, sender, recipient string, stageErr error, stage request.StageResult) error {

	var (
		data,
		dataCrypted,
		reqLinkDecrypted,
		stateBox []byte

		boxName string

		reqLink *request.RequestLink

		err error

		methodName = "SAVE_UPDATE_REQUEST"
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
		return err
	}

	// Десериализуем контейнер
	var boxRequest [][]byte
	err = json.Unmarshal(stateBox, &boxRequest)
	if err != nil {
		bankIDLog.Error(fmt.Sprintf(bankid_constants.ErrorFormatDeserializeContainer, methodName, err))
		return err
	}

	var index = -1
	reqLink = &request.RequestLink{}
	for i, reqLinkCrypted := range boxRequest {

		reqLinkDecrypted, err = oaepHandler.Decrypt(reqLinkCrypted)
		if err != nil {
			bankIDLog.Error(fmt.Sprintf(bankid_constants.ErrorFormatDecryptObjectInternal, methodName, err))
			return err
		}

		err = json.Unmarshal(reqLinkDecrypted, reqLink)
		if err != nil {
			bankIDLog.Error(fmt.Sprintf(bankid_constants.ErrorFormatDeserializeObject, methodName, err))
			return err
		}

		if id == reqLink.Id {
			index = i
			break
		}
	}

	if index == -1 {
		// Создаем линк на запрос
		reqLink = &request.RequestLink{
			Id:        id,
			Sender:    sender,
			Recipient: recipient,
			Stage:     stage,
		}
	} else {
		reqLink.Sender = sender
		reqLink.Recipient = recipient
		reqLink.Stage = stage
	}

	reqLink.SetError(stageErr)

	// Сериализуем линк
	data, err = json.Marshal(reqLink)
	if err != nil {
		bankIDLog.Error(fmt.Sprintf(bankid_constants.ErrorFormatSerializeObject, methodName, err))
		return err
	}

	// Шифруем линк
	dataCrypted, err = oaepHandler.Encrypt(data)
	if err != nil {
		bankIDLog.Error(fmt.Sprintf("[%s]Erorr at encrypt request: %s", methodName, err))
		return err
	}

	// Добавляем или изменеям линк в контейнере
	if index != -1 {
		boxRequest[index] = dataCrypted
	} else {
		boxRequest = append(boxRequest, dataCrypted)
	}

	data, err = json.Marshal(boxRequest)
	if err != nil {
		bankIDLog.Error(fmt.Sprintf(bankid_constants.ErrorFormatDeserializeContainer, methodName, err))
		return err
	}

	// Сохраянем контейнер
	err = stub.PutState(boxName, data)
	if err != nil {
		bankIDLog.Error(fmt.Sprintf(bankid_constants.ErrorFormatSavingState, methodName, err))
		return err
	}

	return nil
}

func stageRequestLink(stub *shim.ChaincodeStub,
	id string, fromStage request.StageResult, toStage request.StageResult) error {

	var (
		data,
		dataCrypted,
		reqLinkDecrypted,
		stateBox []byte

		boxNameCurrent,
		boxNameNext string

		reqLink,
		tmpLink *request.RequestLink

		err error
	)

	// Имя контейнера в котором сейчас находится запрос
	switch fromStage {

	case request.RequestCreate:
		boxNameCurrent = bankid_constants.BoxReqCreatedName

	case request.RequestToUnitValidation:
		boxNameCurrent = bankid_constants.BoxReqToVerificationUnitName

	case request.RequestApproved:
		boxNameCurrent = bankid_constants.BoxReqApprovedName

	case request.RequestRejected:
		boxNameCurrent = bankid_constants.BoxReqRejectedName
	}

	// Имя контейнера в котором будет находится запрос
	switch toStage {

	case request.RequestCreate:
		boxNameNext = bankid_constants.BoxReqCreatedName

	case request.RequestToUnitValidation:
		boxNameNext = bankid_constants.BoxReqToVerificationUnitName

	case request.RequestApproved:
		boxNameNext = bankid_constants.BoxReqApprovedName

	case request.RequestRejected:
		boxNameNext = bankid_constants.BoxReqRejectedName
	}

	// Достаем необходимый контейнер линков
	if stateBox, err = stub.GetState(boxNameCurrent); err != nil {
		return err
	}

	// Десериализуем контейнер
	var currentBoxRequest [][]byte
	if err = json.Unmarshal(stateBox, &currentBoxRequest); err != nil {

		return err
	}

	var currentIndex = -1
	reqLink = &request.RequestLink{}
	for i, reqLinkCrypted := range currentBoxRequest {

		reqLinkDecrypted, err = oaepHandler.Decrypt(reqLinkCrypted)
		if err != nil {

			return err
		}

		err = json.Unmarshal(reqLinkDecrypted, reqLink)
		if err != nil {

			return err
		}

		if id == reqLink.Id {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 {

		// Ошибка, запрос с требуемым статусом не обнаружен в нужном контейнере
		return fmt.Errorf(fmt.Sprintf("Request with id \"%s\" not found in needed box!", id))
	}

	// Ищем запрос в контейнере, в который собираемся переместить запрос
	if stateBox, err = stub.GetState(boxNameNext); err != nil {

		return err
	}

	// Десериализуем контейнер
	var nextBoxRequest [][]byte
	if err = json.Unmarshal(stateBox, &nextBoxRequest); err != nil {

		return err
	}

	var nextIndex = -1
	tmpLink = &request.RequestLink{}
	for i, reqLinkCrypted := range nextBoxRequest {

		reqLinkDecrypted, err = oaepHandler.Decrypt(reqLinkCrypted)
		if err != nil {
			return err
		}

		err = json.Unmarshal(reqLinkDecrypted, tmpLink)
		if err != nil {
			return err
		}

		if id == tmpLink.Id {
			nextIndex = i
			break
		}
	}

	if nextIndex != -1 {

		// Ошибка, запрос уже существует в контейнере, в который мы собираемся переместить его
		return fmt.Errorf(fmt.Sprintf("Request with id \"%s\" already exists in box!", id))
	}

	reqLink.Stage = toStage

	// Сериализуем линк
	data, err = json.Marshal(reqLink)
	if err != nil {

		return err
	}

	// Шифруем линк
	dataCrypted, err = oaepHandler.Encrypt(data)
	if err != nil {

		return err
	}

	// Удаляем реквест из текущего контейнера
	currentBoxRequest = append(currentBoxRequest[:currentIndex], currentBoxRequest[currentIndex+1:]...)

	data, err = json.Marshal(currentBoxRequest)
	if err != nil {

		return err
	}

	// Сохраянем контейнер
	err = stub.PutState(boxNameCurrent, data)
	if err != nil {

		return err
	}

	// Добавляем линк в следующий контейнер
	nextBoxRequest = append(nextBoxRequest, dataCrypted)

	data, err = json.Marshal(nextBoxRequest)
	if err != nil {

		return err
	}

	// Сохраянем контейнер
	err = stub.PutState(boxNameNext, data)
	if err != nil {

		return err
	}

	return nil
}

// Получить запрос из контейнера
func getRequestLink(stub *shim.ChaincodeStub, requestID string, stage request.StageResult) (*request.RequestLink, error) {

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
		bankIDLog.Error(fmt.Sprintf(bankid_constants.ErrorFormatDeserializeContainer, "GET_REQUEST", err))
		return nil, err
	}

	reqLink = &request.RequestLink{}
	for _, reqLinkCrypted = range boxRequest {

		reqLinkDecrypted, err = oaepHandler.Decrypt(reqLinkCrypted)
		if err != nil {
			bankIDLog.Error(fmt.Sprintf(bankid_constants.ErrorFormatDecryptObjectInternal, "GET_REQUEST", err))
			return nil, err
		}

		err = json.Unmarshal(reqLinkDecrypted, reqLink)
		if err != nil {
			bankIDLog.Error(fmt.Sprintf(bankid_constants.ErrorFormatDeserializeObject, "GET_REQUEST", err))
			return nil, err
		}

		if reqLink.Id == requestID {
			return reqLink, nil
		}
	}

	return nil, nil
}

func checkUnitExist(stub *shim.ChaincodeStub, unitAlias string) bool {

	unit, err := stub.GetState(unitAlias)
	if err != nil || unit == nil {

		return false
	}

	return true
}

func init() {

	bankIDLog = logging.MustGetLogger("chaincode/bankid")

	logFormatter := logging.NewBackendFormatter(
		logging.NewLogBackend(os.Stderr, "", 0),
		logging.MustStringFormatter(
			`%{color}%{time:2006/01/02 15:04:05}%{color:reset} %{message}`,
		))
	logging.SetBackend(logFormatter)

	// Инициализируем oaep
	initOaep()

	// Инициализируем валидаторы
	initValidaters()

	// Инициализируем обработчики запросов
	initChaincodeHandlers()
}

func initOaep() {

	oaepHandler = &oaep.OaepHandler{}

	oaepHandler.LoadPrivateKeyFile("id_rsa")
	oaepHandler.LoadPublicKey(genesis.PubKeyBankId)

	if !oaepHandler.Initialized() {
		panic("OaepHadler is not initialized! Terminated...")
	}
}

func initValidaters() {

	validatorRequestIdentifier = &validators.ValidatorRequestIdentifier{
		Log: bankIDLog,
	}

	parameterValidaters = append(parameterValidaters, validatorRequestIdentifier)
}

func initChaincodeHandlers() {

	queryHandler = &chaincode_call_handlers.QueryHandler{}
	queryHandler.Init(bankIDLog, oaepHandler)
}

func main() {

	err := shim.Start(new(BankIDChainCode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

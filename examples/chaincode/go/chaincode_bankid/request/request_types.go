package request

import (
	"time"
	"fmt"
)

// Запрос на подверждение личности клиента
type RequestPersonValidation struct {

	// Идентификатор запроса
	Id string

	// Этап выполнения запроса
	Stage *Stage

	// Данные запроса
	Person *Person

	// Центр где был инициирован запрос
	Sender string

	// Центр в котором будет проверен клиент
	Recipient string

	// Тип идентификации клиента
	TypeIdentification TypeIdentification

	// Дата создания
	Created time.Time

	// Дата последнего изменения
	LastModified time.Time

	// Результат выполнения идентификации
	Result interface{}
}

// Ошибка выполнения идентификации
type ErrorRequest interface {

	Error() error

	SetError(err error)
}

// Ссылка на запрос
type RequestLink struct {

	Id 			string

	Sender 		string

	Recipient 	string

	Stage 		StageResult

	err 		error
}

// Информация о клиенте
type Person struct {

	// Идентификатор клиента в системе, в которой он авторизован
	ClientId 				string

	// Тип набора данных для идентификации личности
	TypeSetIdentification 	TypeSetIdentification

	// Хеш персональных данных
	HashPersonalInfo 		string
}

// Тип идентификации клиента
type TypeSetIdentification int
const (

	// Стандартный набор для идентфиикации клиента.
	// Хеш производится на основе ПД клиента: HASH(Фамилия;Имя;Отчество;Серия паспорта;Номер паспорта)
	TypeSetIdentificationStandard = 1 + iota

	// Банковский набор для идентфиикации клиента.
	// Хеш производится на основе расчетного счета клиента: HASH(Расчетный счет)
	TypeSetIdentificationBankSet
)

// Этап в котором находится запрос на валидацию клиента
type Stage struct {

	// Результат выполнения этапа
	StageResult 	StageResult

	// Ошибка выполнения этапа если она есть
	Error 		error
}

// Результат этапа выполнения запроса на идентфикацию
type StageResult int
const (

	// Запрос создан с ошибкой
	RequestCreate = 1 + iota

	// Запрос проверен и отправлен на идентификацию в recipient-центр
	RequestToUnitValidation

	// Запрос проверен и утвержден
	RequestApproved

	// Запрос на идентфикацию клиента отклонен
	RequestRejected
)

// Тип идентификации клиента
type TypeIdentification int
const (

	// Клиент уже был идентифицирован одним из центров
	TypeIdentificationClientIdentified = 1 + iota

	// Клиент еще не был идентифицирован
	TypeIdentificationClientNotIdentified
)

// Интерфейс получения результата выполнения запроса на идентификацию
type RequestResultManager interface {

	SetResult(data interface{})

	GetResult() interface{}
}

type ApprovedData struct {

	Url string
}

type RejectedData struct {

	Reason string
}

func (r *RequestLink) Error() error {
	return r.err
}

func (r *RequestLink) SetError(err error) {
	r.err = err
}

func (r *RequestPersonValidation) SetResult(data interface{}) {

	r.Result = data
}

func (r *RequestPersonValidation) GetResult() interface{} {

	return r.Result
}

func (r *RequestPersonValidation) String() string {

	var stageDescr string
	switch r.Stage.StageResult {

	case RequestCreate:
		stageDescr = "Request created with error"

	case RequestToUnitValidation:
		stageDescr = "Request await unit validation"

	case RequestApproved:
		stageDescr = "Request was approved"

	case RequestRejected:
		stageDescr = "Request was rejected"
	}

	var typeIdentificationDescr string
	switch r.TypeIdentification {

	case TypeIdentificationClientIdentified:
		typeIdentificationDescr = "PRE-IDENTIFICATION"

	case TypeIdentificationClientNotIdentified:
		typeIdentificationDescr = "POST-IDENTIFICATION"
	}

	var typeSetidentification string
	switch r.Person.TypeSetIdentification {

	case TypeSetIdentificationStandard:
		typeSetidentification = "STANDART"

	case TypeSetIdentificationBankSet:
		typeSetidentification = "BANKSET"
	}

	var result = ""
	if (r.Stage.StageResult == RequestApproved ||
		r.Stage.StageResult == RequestRejected) {

		resultObj  := r.GetResult()
		if resultObj != nil {

			switch r.Stage.StageResult {

			case RequestApproved:
				result = resultObj.(map[string]interface{})["Url"].(string)

			case RequestRejected:
				result = resultObj.(map[string]interface{})["Reason"].(string)
			}
		}

	}

	return fmt.Sprintf("{ Id: \"%s\", Stage: \"%s\", Error: %#v, Sender: \"%s\", Recipient: \"%s\", TypeIdentification: \"%s\", " +
		"Created: \"%s\", LastModified: \"%s\", Person: [%s] }, Result: \"%s\"",
		r.Id,
		stageDescr,
		r.Stage.Error,
		r.Sender,
		r.Recipient,
		typeIdentificationDescr,
		r.Created,
		r.LastModified,
		fmt.Sprintf("Id - %s, Hash - %s, Set - %s",
			r.Person.ClientId, r.Person.HashPersonalInfo, typeSetidentification),
		result,
	)
}
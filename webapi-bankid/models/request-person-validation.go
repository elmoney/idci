package models

import "time"

//RequestPersonValidation Запрос на подверждение личности клиента
type RequestPersonValidation struct {

	// Идентификатор запроса
	ID                 string

	// Этап выполнения запроса
	Stage              *Stage

	// Данные запроса
	Person             *Person

	// Центр где был инициирован запрос
	Sender             string

	// Центр в котором будет проверен клиент
	Recipient          string

	// Тип идентификации клиента
	TypeIdentification TypeIdentification

	// Дата создания
	Created            time.Time

	// Дата последнего изменения
	LastModified       time.Time

	// Результат выполнения идентификации
	Result             interface{}
}

//ErrorRequest Ошибка выполнения идентификации
type ErrorRequest interface {

	Error() error

	SetError(err error)
}

//RequestLink Ссылка на запрос
type RequestLink struct {

	ID        string

	Sender    string

	Recipient string

	Stage     StageResult

	err       error
}

//Person Информация о клиенте
type Person struct {

	// Идентификатор клиента в системе, в которой он авторизован
	ClientID              string

	// Тип набора данных для идентификации личности
	TypeSetIdentification TypeSetIdentification

	// Хеш персональных данных
	HashPersonalInfo      string
}

//TypeSetIdentification Тип идентификации клиента
type TypeSetIdentification int
const (

//TypeSetIdentificationStandard Стандартный набор для идентфиикации клиента.
// Хеш производится на основе расчетного счета клиента: HASH(Фамилия;Имя;Отчество;Серия паспорта;Номер паспорта)
	TypeSetIdentificationStandard = 1 + iota

//TypeSetIdentificationBankSet Банковский набор для идентфиикации клиента.
// Хеш производится на основе расчетного счета клиента: HASH(Расчетный счет)
	TypeSetIdentificationBankSet
)

//Stage  Этап в котором находится запрос на валидацию клиента
type Stage struct {

	// Результат выполнения этапа
	StageResult StageResult

	// Ошибка выполнения этапа если она есть
	Error 		error
}

//StageResult Результат этапа выполнения запроса на идентфикацию
type StageResult int
const (

//RequestCreate Запрос создан с ошибкой
	RequestCreate = 1 + iota

//RequestToUnitValidation Запрос проверен и отправлен на идентификацию в recipient-центр
	RequestToUnitValidation

//RequestApproved Запрос проверен и утвержден
	RequestApproved

//RequestRejected Запрос на идентфикацию клиента отклонен
	RequestRejected
)

//TypeIdentification Тип идентификации клиента
type TypeIdentification int
const (

//TypeIdentificationClientIdentified Клиент уже был идентифицирован одним из центров
	TypeIdentificationClientIdentified = 1 + iota

//TypeIdentificationClientNotIdentified Клиент еще не был идентифицирован
	TypeIdentificationClientNotIdentified
)

//RequestResultManager Интерфейс получения результата выполнения запроса на идентификацию
type RequestResultManager interface {

	SetResult(data interface{})

	GetResult() interface{}
}

//ApprovedData блок для хранения данных по идентифкации
type ApprovedData struct {

	URL string
}

//RejectedData блок для хранения данных по идентифкации
type RejectedData struct {

	Reason string
}
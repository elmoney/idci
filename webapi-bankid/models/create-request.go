package models


//CreateRequest модель для создания запроса на идентифкацию
type CreateRequest struct {
	RecipientAlias     string `json:"recipientAlias"`
	ClientID           string `json:"clientId"`
	TypeIdentification string `json:"typeIdentification"`
	IdentificationSet  struct {
		SetType      int `json:"setType"`
		PersonalData struct {
			FirstName  string `json:"firstName"`
			LastName   string `json:"lastName"`
			MiddleName string `json:"middleName"`
			PasSer     string `json:"pasSer"`
			PasNumber  string `json:"pasNumber"`
		} `json:"personalData"`
		BankSet struct {
			AccountNumber string `json:"accountNumber"`
		} `json:"bankSet"`
	} `json:"identificationSet"`
}

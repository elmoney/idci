package models


//UnitIdentifier  тип идентфикатора в системе участника системы
type  UnitIdentifier struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Regexp      string  `json:"regexp"`
}

package models



//Unit описание  участника системы
type Unit struct {
	Alias      string `json:"alias"`
	FullName   string `json:"fullname"`
	Status     int `json:"status"`
	PublicKey  []byte `json:"publicKey"`
	Identifier UnitIdentifier `json:"identifier"`
	Partners   []string `json:"partners"`
}


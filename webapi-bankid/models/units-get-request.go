package models


//UnitsGetRequest запрос на получение данных по учасникам системы
type UnitsGetRequest struct {
	OK struct {
		   Units []Unit `json:"units"`
	   } `json:"OK"`
	Error string `json:"Error"`
}

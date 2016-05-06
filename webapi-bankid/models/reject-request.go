package models

//RejectRequest  модель для отклонения запроса на идентификацию
type RejectRequest struct {
	RequestID string `json:"requestId"`
	Reason    string `json:"reason"`
}

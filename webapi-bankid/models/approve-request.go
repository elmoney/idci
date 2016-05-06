package models

//ApproveRequest модель для  подтвержденя запроса на идентификацию
type ApproveRequest struct {
	RequestID string `json:"requestId"`
	URL       string `json:"url"`
}

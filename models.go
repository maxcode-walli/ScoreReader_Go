package p

type TransactionScore struct {
	UserID            string `json:"userId,omitempty"`
	ExternalAccountId string `json:"externalAccountID,omitempty"`
	TransactionID     string `json:"transactionID,omitempty"`
	Label             string `json:"label,omitempty"`
	Score             string `json:"score,omitempty"`
}

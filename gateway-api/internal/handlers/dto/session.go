package dto

type SessionResponse struct {
	ID           string               `json:"id"`
	Name         string               `json:"name"`
	Total        float64              `json:"total"`
	IsActive     bool                 `json:"isActive"`
	Members      []SessionMember      `json:"members"`
	Transactions []SessionTransaction `json:"transactions"`
}

type SessionMember struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type SessionTransactionMember struct {
	ID       string  `json:"userId"`
	Name     string  `json:"name"`
	Username string  `json:"username"`
	Amount   float64 `json:"amount"`
}

type SessionTransaction struct {
	ID        string                     `json:"id"`
	CreatorID string                     `json:"creatorId"`
	Total     float64                    `json:"total"`
	Members   []SessionTransactionMember `json:"members"`
}

package events

type UserCreatedEventData struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
}

type UserCreatedEvent struct {
	User UserCreatedEventData `json:"user"`
}

type SessionTransactionCreatedEvent struct {
	SessionID     string          `json:"session_id"`
	TransactionID string          `json:"transaction_id"`
	CreatorID     string          `json:"creator_id"`
	Total         float64         `json:"total"`
	Members       []SessionMember `json:"members"`
}

type SessionMember struct {
	ID string `json:"id"`
}

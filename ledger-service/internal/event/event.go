package event

type TransactionCreatedMemberAmount struct {
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
}

// TransactionCreatedEvent is an event that is published when a transaction is created.
// The CreatorID indicates the user who created the transaction.
// The MemberAmounts slice contains the user IDs and amounts for each member of the transaction.
type TransactionCreatedEvent struct {
	TransactionID string                           `json:"transaction_id"`
	CreatorID     string                           `json:"creator_id"`
	SessionID     string                           `json:"session_id"`
	MemberAmounts []TransactionCreatedMemberAmount `json:"member_amounts"`
}

// LedgerUpdateEvent is an event that is published when the ledger is updated.
// This event contains a single item from the ledger.
// The transaction will be between the creator and one of the members or one of the members and the creator.
type LedgerUpdateEvent struct {
	ID            string  `json:"id"`
	TransactionID string  `json:"transaction_id"`
	SessionID     string  `json:"session_id"`
	UserID        string  `json:"user_id"`
	ParticipantID string  `json:"participant_id"`
	Amount        float64 `json:"amount"`
}

// UserTotalsEvent describes the total amounts of a specific user at the time the event was sent.
type UserTotalsEvent struct {
	UserID string  `json:"user_id"`
	Credit float64 `json:"credit"`
	Debit  float64 `json:"debit"`
	Net    float64 `json:"net"`
}

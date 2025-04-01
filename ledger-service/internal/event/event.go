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
// This event contains a single transaction from the ledger.
// The transaction will be between the creator and one of the members or one of the members and the creator.
type LedgerUpdateEvent struct {
	ID            string  `json:"id"`
	TransactionID string  `json:"transaction_id"`
	SessionID     string  `json:"session_id"`
	UserID        string  `json:"user_id"`
	ParticipantID string  `json:"participant_id"`
	Amount        float64 `json:"amount"`
}

type LedgerUpdateCompleteMemberAmount struct {
	UserID string  `json:"user_id"`
	Amount float64 `json:"amount"`
}

type LedgerTransactionUpdatedEvent struct {
	TransactionID string                             `json:"transaction_id"`
	SessionID     string                             `json:"session_id"`
	UserID        string                             `json:"user_id"`
	Total         float64                            `json:"total"`
	Amounts       []LedgerUpdateCompleteMemberAmount `json:"amounts"`
}

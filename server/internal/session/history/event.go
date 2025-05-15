package history

type EventType string

const (
	EventUnknown                EventType = "unknown_event"
	EventTransactionCreated     EventType = "transaction_created"
	EventMemberAdded            EventType = "member_added"
	EventMemberRemoved          EventType = "member_removed"
	EventMemberLeft             EventType = "member_left"
	EventSessionClosed          EventType = "session_closed"
	EventSessionOpened          EventType = "session_opened"
	EventMemberPromotedToAdmin  EventType = "promoted_to_admin"
	EventMemberDemotedFromAdmin EventType = "demoted_from_admin"
)

func NewEventType(t string) EventType {
	switch t {
	case EventTransactionCreated.String():
		return EventTransactionCreated
	case EventMemberAdded.String():
		return EventMemberAdded
	case EventMemberRemoved.String():
		return EventMemberRemoved
	case EventMemberLeft.String():
		return EventMemberLeft
	case EventSessionClosed.String():
		return EventSessionClosed
	case EventSessionOpened.String():
		return EventSessionOpened
	case EventMemberPromotedToAdmin.String():
		return EventMemberPromotedToAdmin
	case EventMemberDemotedFromAdmin.String():
		return EventMemberDemotedFromAdmin
	default:
		return EventUnknown
	}
}

func (et EventType) String() string {
	return string(et)
}

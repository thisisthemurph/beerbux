package events

type UserCreatedEventData struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
}

type UserCreatedEvent struct {
	User UserCreatedEventData `json:"user"`
}

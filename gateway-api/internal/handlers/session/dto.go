package session

type SessionResponse struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	IsActive bool            `json:"isActive"`
	Members  []SessionMember `json:"members"`
}

type SessionMember struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

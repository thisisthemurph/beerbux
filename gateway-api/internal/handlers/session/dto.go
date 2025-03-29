package session

type SessionResponse struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Total    float64         `json:"total"`
	IsActive bool            `json:"isActive"`
	Members  []SessionMember `json:"members"`
}

type SessionMember struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

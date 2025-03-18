package user

type UserResponse struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Username   string  `json:"username"`
	NetBalance float64 `json:"netBalance"`
}

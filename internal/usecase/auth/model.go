package auth

type UserView struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
}

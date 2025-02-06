package models

type DeskyUser struct {
	ID       uint   `json:"id,omitempty"`
	Login    string `json:"login"`
	Password string `json:"password,omitempty"`
}

func NewDeskyUser(id uint, login, pwd string) *DeskyUser {
	return &DeskyUser{
		ID:       id,
		Login:    login,
		Password: pwd,
	}
}

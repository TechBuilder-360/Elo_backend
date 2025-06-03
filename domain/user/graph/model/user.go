package model

type Login struct {
	Email string `json:"email"`
}

type LoginResult struct {
	AuthenticationCode string `json:"authentication_code"`
}

package user

type AuthenticateParam struct {
	Email    string
	Password string
}

type AuthenticateResult struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

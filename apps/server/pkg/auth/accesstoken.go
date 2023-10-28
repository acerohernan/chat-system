package auth

type AccessToken struct {
	Grants *Grants
	Jwt    string
}

type Grants struct {
	Email          string `json:"email"`
	CanSendMessage bool   `json:"canSendMessage"`
}

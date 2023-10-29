package auth

type AccessToken struct {
	Grants *Grants
	Jwt    string
}

type Grants struct {
	Id             string `json:"id"`
	Email          string `json:"email"`
	CanSendMessage bool   `json:"canSendMessage"`
}

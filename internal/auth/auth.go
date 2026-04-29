package auth

type Claims struct {
	Subject  string
	Username string
}

type Authenticator interface {
	ValidateToken(token string) (*Claims, error)
	GenerateToken(userID int64, username string) (string, error)
}

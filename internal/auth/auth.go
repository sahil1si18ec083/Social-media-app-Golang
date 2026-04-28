package auth

type Authenticator interface {
	ValidateToken(token string) error
	GenerateToken(userID int64, username string) (string, error)
}

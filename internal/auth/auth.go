package auth

type Authenticator interface {
	ValidateToken(token string) error
}

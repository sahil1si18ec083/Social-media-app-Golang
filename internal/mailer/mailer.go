package mailer

const (
	FromName            = "GopherSocial"
	maxRetires          = 3
	UserWelcomeTemplate = "user_invitation.tmpl"
)

type Client interface {
	Send(templateFile, username, email string, data any, isSandbox bool) (int, error)
}

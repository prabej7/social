package mailer

import "embed"

var (
	FromName            = "LopherSocial"
	MaxRetries          = 3
	UserWelcomeTemplate = "user_invitation.templ"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile, username, email string, data any, isSandbox bool) error
}

package mailer

import "embed"

//go:embed template/*
var FS embed.FS

const ( 
	FromName = "GopherSocial"
	maxRetries = 3
	UserWelcomeTemplate = "user_invitation.tmpl"
)

type Client interface{
	Send(templateFile, username,email string, data any, isDevelopment bool) (int,error)
}


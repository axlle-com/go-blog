package contracts

type MailRequest interface {
	From() string
	To() string
	Subject() string
	Body() string
	ToString() string
}

type Mailer interface {
	SendMail(MailRequest)
}

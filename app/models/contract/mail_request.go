package contract

type MailRequest interface {
	GetFrom() string
	GetTo() string
	GetSubject() string
	GetBody() string
	ToString() string
}

type Mailer interface {
	SendMail(MailRequest) error
}

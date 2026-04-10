package email

type SendEmailRes struct {
	EmailId string
}

type SendEmailReq struct {
	Subject string

	ToEmail string
	ToName  string

	BodyHtml string
	BodyText string
}

package mailpit

// SendEmailRequest represents the Mailpit API request body for sending an email.
type SendEmailRequest struct {
	Attachments []Attachment      `json:"Attachments,omitzero"`
	Bcc         []string          `json:"Bcc,omitzero"`
	Cc          []Recipient       `json:"Cc,omitzero"`
	From        Recipient         `json:"From,omitzero"`
	HTML        string            `json:"HTML,omitzero"`
	Headers     map[string]string `json:"Headers,omitzero"`
	ReplyTo     []Recipient       `json:"ReplyTo,omitzero"`
	Subject     string            `json:"Subject,omitzero"`
	Tags        []string          `json:"Tags,omitzero"`
	Text        string            `json:"Text,omitzero"`
	To          []Recipient       `json:"To,omitzero"`
}

// Attachment represents an email attachment in the Mailpit request.
type Attachment struct {
	Content     string `json:"Content"`
	ContentID   string `json:"ContentID,omitzero"`
	ContentType string `json:"ContentType"`
	Filename    string `json:"Filename"`
}

// Recipient represents an email recipient (Name and Email).
type Recipient struct {
	Email string `json:"Email"`
	Name  string `json:"Name,omitzero"`
}

type SendEmailResponse struct {
	ID string `json:"ID"`
}

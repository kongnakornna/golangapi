package presenter

// EmailSendRequest represents an email sending request.
// คำขอส่งอีเมล
type EmailSendRequest struct {
	To      string `json:"to" validate:"required,email"`
	Cc      string `json:"cc"`
	Bcc     string `json:"bcc"`
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body" validate:"required"`
}

// EmailLogResponse represents an email log entry.
// บันทึกการส่งอีเมล
type EmailLogResponse struct {
	ID           int    `json:"id"`
	To           string `json:"to"`
	Subject      string `json:"subject"`
	Status       string `json:"status"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	SentAt       string `json:"sentAt,omitempty"`
	CreatedAt    string `json:"createdAt"`
}

// EmailConfigResponse represents the SMTP configuration.
// การตั้งค่า SMTP
type EmailConfigResponse struct {
	SmtpHost  string `json:"smtpHost"`
	SmtpPort  int    `json:"smtpPort"`
	SmtpUser  string `json:"smtpUser"`
	FromEmail string `json:"fromEmail"`
	FromName  string `json:"fromName"`
	IsActive  bool   `json:"isActive"`
}

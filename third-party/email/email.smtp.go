package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"gopkg.in/gomail.v2"
)

// SMTP Sender Implementation
func NewSMTPSender(cfg Config) EmailSender {
	return &SMTPSender{Config: cfg}
}

func (s *SMTPSender) GetProvider() Provider {
	return ProviderSMTP
}

func (s *SMTPSender) SendEmail(ctx context.Context, email *Email) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.Config.From)
	m.SetHeader("To", email.To)
	m.SetHeader("Subject", email.Subject)

	if email.ReplyTo != "" {
		m.SetHeader("Reply-To", email.ReplyTo)
	}

	if len(email.CC) > 0 {
		m.SetHeader("Cc", email.CC...)
	}

	if len(email.BCC) > 0 {
		m.SetHeader("Bcc", email.BCC...)
	}

	// Handle body content
	if err := s.setEmailBody(email, m); err != nil {
		return err
	}

	d := gomail.NewDialer(s.Config.Host, s.Config.Port, s.Config.Username, s.Config.Password)
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("SMTP failed to send email: %w", err)
	}

	return nil
}

func (s *SMTPSender) setEmailBody(email *Email, m *gomail.Message) error {
	if email.TemplatePath != "" && email.TemplateData != nil {
		tmpl, err := template.ParseFiles(email.TemplatePath)
		if err != nil {
			return fmt.Errorf("failed to parse template: %w", err)
		}

		var body bytes.Buffer
		if err := tmpl.Execute(&body, email.TemplateData); err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}

		m.SetBody("text/html", body.String())
	} else if email.BodyHTML != "" {
		if email.BodyText != "" {
			m.SetBody("text/plain", email.BodyText)
			m.AddAlternative("text/html", email.BodyHTML)
		} else {
			m.SetBody("text/html", email.BodyHTML)
		}
	} else if email.BodyText != "" {
		m.SetBody("text/plain", email.BodyText)
	} else {
		return fmt.Errorf("no email body provided")
	}

	return nil
}

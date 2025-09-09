package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"github.com/resend/resend-go/v2"
)

// RESEND EMAIl

// Resend Sender Implementation
func NewResendSender(cfg Config) EmailSender {
	client := resend.NewClient(cfg.APIKey)
	return &ResendSender{
		Config: cfg,
		Client: client,
	}
}

func (r *ResendSender) GetProvider() Provider {
	return ProviderResend
}

func (r *ResendSender) SendEmail(ctx context.Context, email *Email) error {
	// Handle template rendering if needed
	htmlBody, textBody, err := r.prepareEmailBody(email)
	if err != nil {
		return err
	}

	params := &resend.SendEmailRequest{
		From:    r.Config.From,
		To:      []string{email.To},
		Subject: email.Subject,
	}

	if len(email.CC) > 0 {
		params.Cc = email.CC
	}

	if len(email.BCC) > 0 {
		params.Bcc = email.BCC
	}

	if email.ReplyTo != "" {
		params.ReplyTo = email.ReplyTo
	}

	if htmlBody != "" {
		params.Html = htmlBody
	}

	if textBody != "" {
		params.Text = textBody
	}

	// Add tags if supported
	if len(email.Tags) > 0 {
		params.Tags = []resend.Tag{}
		for key, value := range email.Tags {
			params.Tags = append(params.Tags, resend.Tag{
				Name:  key,
				Value: value,
			})
		}
	}

	_, err = r.Client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("Resend failed to send email: %w", err)
	}

	return nil
}

func (r *ResendSender) prepareEmailBody(email *Email) (string, string, error) {
	var htmlBody, textBody string

	if email.TemplatePath != "" && email.TemplateData != nil {
		tmpl, err := template.ParseFiles(email.TemplatePath)
		if err != nil {
			return "", "", fmt.Errorf("failed to parse template: %w", err)
		}

		var body bytes.Buffer
		if err := tmpl.Execute(&body, email.TemplateData); err != nil {
			return "", "", fmt.Errorf("failed to execute template: %w", err)
		}

		htmlBody = body.String()
	} else {
		htmlBody = email.BodyHTML
		textBody = email.BodyText
	}

	if htmlBody == "" && textBody == "" {
		return "", "", fmt.Errorf("no email body provided")
	}

	return htmlBody, textBody, nil
}

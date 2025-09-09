package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

// NewSESSender SES Sender Implementation
func NewSESSender(cfg Config) (EmailSender, error) {
	// Note: You'll need to properly configure AWS SDK with credentials
	// This is a simplified example
	if cfg.Region == "" {
		return nil, fmt.Errorf("AWS region is required for SES")
	}

	// Initialize AWS SES client here
	// client := ses.NewFromConfig(awsConfig)

	return &SESSender{
		Config: cfg,
		// Client: client,
	}, nil
}

func (s *SESSender) GetProvider() Provider {
	return ProviderSES
}

func (s *SESSender) SendEmail(ctx context.Context, email *Email) error {
	// Handle template rendering if needed
	htmlBody, textBody, err := s.prepareEmailBody(email)
	if err != nil {
		return err
	}

	destinations := []string{email.To}
	destinations = append(destinations, email.CC...)
	destinations = append(destinations, email.BCC...)

	input := &ses.SendEmailInput{
		Source: aws.String(s.Config.From),
		Destination: &types.Destination{
			ToAddresses:  []string{email.To},
			CcAddresses:  email.CC,
			BccAddresses: email.BCC,
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data: aws.String(email.Subject),
			},
			Body: &types.Body{},
		},
	}

	if email.ReplyTo != "" {
		input.ReplyToAddresses = []string{email.ReplyTo}
	}

	if htmlBody != "" {
		input.Message.Body.Html = &types.Content{
			Data: aws.String(htmlBody),
		}
	}

	if textBody != "" {
		input.Message.Body.Text = &types.Content{
			Data: aws.String(textBody),
		}
	}

	// Add tags if supported
	if len(email.Tags) > 0 {
		var tags []types.MessageTag
		for key, value := range email.Tags {
			tags = append(tags, types.MessageTag{
				Name:  aws.String(key),
				Value: aws.String(value),
			})
		}
		input.Tags = tags
	}

	_, err = s.Client.SendEmail(ctx, input)
	if err != nil {
		return fmt.Errorf("SES failed to send email: %w", err)
	}

	return nil
}

func (s *SESSender) prepareEmailBody(email *Email) (string, string, error) {
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

package email

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/resend/resend-go/v2"
)

// Provider represents different email service providers
type Provider string

const (
	ProviderSMTP   Provider = "smtp"
	ProviderSES    Provider = "ses"
	ProviderResend Provider = "resend"
)

// Config holds configuration for different email providers
type Config struct {
	Provider Provider
	// SMTP specific
	Host     string
	Port     int
	Username string
	Password string
	From     string

	// AWS SES specific
	Region    string
	AccessKey string
	SecretKey string

	// Resend specific
	APIKey string
}

// EmailSender defines the interface for sending emails
type EmailSender interface {
	SendEmail(ctx context.Context, email *Email) error
	GetProvider() Provider
}

// EmailBuilder defines the interface for building emails
type EmailBuilder interface {
	To(to string) EmailBuilder
	Subject(subject string) EmailBuilder
	BodyHTML(body string) EmailBuilder
	BodyText(body string) EmailBuilder
	BodyFromTemplate(templatePath string, data interface{}) EmailBuilder
	CC(cc ...string) EmailBuilder
	BCC(bcc ...string) EmailBuilder
	ReplyTo(replyTo string) EmailBuilder
	Tag(key, value string) EmailBuilder
	Build() (*Email, error)
	Send(ctx context.Context) error
}

// Email represents an email message
type Email struct {
	To           string
	CC           []string
	BCC          []string
	Subject      string
	BodyHTML     string
	BodyText     string
	ReplyTo      string
	TemplatePath string
	TemplateData interface{}
	Tags         map[string]string
}

// SMTP implementation
type SMTPSender struct {
	Config Config
}

// SESSender SES implementation
type SESSender struct {
	Config Config
	Client *ses.Client
}

// ResendSender Resend implementation
type ResendSender struct {
	Config Config
	Client *resend.Client
}

// emailBuilder implements EmailBuilder interface
type emailBuilder struct {
	email  *Email
	sender EmailSender
	errors []error
}

// NewEmailSender Factory function to create appropriate sender
func NewEmailSender(cfg Config) (EmailSender, error) {
	switch cfg.Provider {
	case ProviderSMTP:
		return NewSMTPSender(cfg), nil
	case ProviderSES:
		return NewSESSender(cfg)
	case ProviderResend:
		return NewResendSender(cfg), nil
	default:
		return nil, fmt.Errorf("unsupported email provider: %s", cfg.Provider)
	}
}

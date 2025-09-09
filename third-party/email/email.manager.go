package email

import (
	"context"
	"fmt"
	"os"
	"strconv"
)

// EmailManager provides a high-level interface for managing multiple email providers
type EmailManager struct {
	providers map[Provider]EmailSender
	primary   Provider
	fallback  Provider
}

// ProviderConfig holds provider-specific configuration
type ProviderConfig struct {
	SMTP   *Config
	SES    *Config
	Resend *Config
}

// NewEmailManager creates a new email manager with multiple providers
func NewEmailManager(configs ProviderConfig) (*EmailManager, error) {
	manager := &EmailManager{
		providers: make(map[Provider]EmailSender),
	}

	// Initialize SMTP if configured
	if configs.SMTP != nil {
		configs.SMTP.Provider = ProviderSMTP
		sender, err := NewEmailSender(*configs.SMTP)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize SMTP: %w", err)
		}
		manager.providers[ProviderSMTP] = sender
		if manager.primary == "" {
			manager.primary = ProviderSMTP
		}
	}

	// Initialize SES if configured
	if configs.SES != nil {
		configs.SES.Provider = ProviderSES
		sender, err := NewEmailSender(*configs.SES)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize SES: %w", err)
		}
		manager.providers[ProviderSES] = sender
		if manager.primary == "" {
			manager.primary = ProviderSES
		}
	}

	// Initialize Resend if configured
	if configs.Resend != nil {
		configs.Resend.Provider = ProviderResend
		sender, err := NewEmailSender(*configs.Resend)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize Resend: %w", err)
		}
		manager.providers[ProviderResend] = sender
		if manager.primary == "" {
			manager.primary = ProviderResend
		}
	}

	if len(manager.providers) == 0 {
		return nil, fmt.Errorf("at least one email provider must be configured")
	}

	return manager, nil
}

// SetPrimaryProvider sets the primary email provider
func (em *EmailManager) SetPrimaryProvider(provider Provider) error {
	if _, exists := em.providers[provider]; !exists {
		return fmt.Errorf("provider %s is not configured", provider)
	}
	em.primary = provider
	return nil
}

// SetFallbackProvider sets the fallback email provider
func (em *EmailManager) SetFallbackProvider(provider Provider) error {
	if _, exists := em.providers[provider]; !exists {
		return fmt.Errorf("provider %s is not configured", provider)
	}
	em.fallback = provider
	return nil
}

// GetProvider returns a specific email provider
func (em *EmailManager) GetProvider(provider Provider) (EmailSender, error) {
	sender, exists := em.providers[provider]
	if !exists {
		return nil, fmt.Errorf("provider %s is not configured", provider)
	}
	return sender, nil
}

// GetPrimaryProvider returns the primary email provider
func (em *EmailManager) GetPrimaryProvider() EmailSender {
	return em.providers[em.primary]
}

// NewBuilder creates a new email builder using the primary provider
func (em *EmailManager) NewBuilder() EmailBuilder {
	return NewEmailBuilder(em.providers[em.primary])
}

// NewBuilderWithProvider creates a new email builder using a specific provider
func (em *EmailManager) NewBuilderWithProvider(provider Provider) (EmailBuilder, error) {
	sender, exists := em.providers[provider]
	if !exists {
		return nil, fmt.Errorf("provider %s is not configured", provider)
	}
	return NewEmailBuilder(sender), nil
}

// SendWithFallback attempts to send email with primary provider, falls back to fallback provider on failure
func (em *EmailManager) SendWithFallback(ctx context.Context, email *Email) error {
	// Try primary provider
	if err := em.providers[em.primary].SendEmail(ctx, email); err != nil {
		if em.fallback != "" && em.fallback != em.primary {
			// Try fallback provider
			if fallbackErr := em.providers[em.fallback].SendEmail(ctx, email); fallbackErr != nil {
				return fmt.Errorf("primary provider (%s) failed: %w, fallback provider (%s) failed: %w",
					em.primary, err, em.fallback, fallbackErr)
			}
			return nil // Fallback succeeded
		}
		return fmt.Errorf("primary provider (%s) failed and no fallback configured: %w", em.primary, err)
	}
	return nil // Primary succeeded
}

// GetAvailableProviders returns a list of configured providers
func (em *EmailManager) GetAvailableProviders() []Provider {
	var providers []Provider
	for provider := range em.providers {
		providers = append(providers, provider)
	}
	return providers
}

func ConfigFromEnv() ProviderConfig {
	config := ProviderConfig{}

	if host := os.Getenv("SMTP_HOST"); host != "" {
		port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
		if port == 0 {
			port = 587 // Default SMTP port
		}
		config.SMTP = &Config{
			Host:     host,
			Port:     port,
			Username: os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("SMTP_PASSWORD"),
			From:     os.Getenv("SMTP_FROM"),
		}
	}

	// AWS SES Configuration
	if region := os.Getenv("AWS_REGION"); region != "" {
		config.SES = &Config{
			Region:    region,
			AccessKey: os.Getenv("AWS_ACCESS_KEY_ID"),
			SecretKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
			From:      os.Getenv("SES_FROM_EMAIL"),
		}
	}

	// Resend Configuration
	if apiKey := os.Getenv("RESEND_API_KEY"); apiKey != "" {
		config.Resend = &Config{
			APIKey: apiKey,
			From:   os.Getenv("RESEND_FROM_EMAIL"),
		}
	}

	return config
}

// EmailService provides a high-level service interface
type EmailService struct {
	manager *EmailManager
}

// NewEmailService creates a new email service
func NewEmailService(configs ProviderConfig) (*EmailService, error) {
	manager, err := NewEmailManager(configs)
	if err != nil {
		return nil, err
	}

	return &EmailService{manager: manager}, nil
}

// NewEmailServiceFromEnv creates a new email service from environment variables
func NewEmailServiceFromEnv() (*EmailService, error) {
	configs := ConfigFromEnv()
	return NewEmailService(configs)
}

// SendWelcomeEmail sends a welcome email using a template
func (es *EmailService) SendWelcomeEmail(ctx context.Context, to, name string) error {
	data := struct {
		Name string
		Year int
	}{
		Name: name,
		Year: 2024,
	}

	return es.manager.NewBuilder().
		To(to).
		Subject("Welcome to Our Platform!").
		BodyFromTemplate("templates/welcome.html", data).
		Tag("type", "welcome").
		Tag("automated", "true").
		Send(ctx)
}

// SendPasswordResetEmail sends a password reset email
func (es *EmailService) SendPasswordResetEmail(ctx context.Context, to, resetLink string) error {
	data := struct {
		ResetLink   string
		ExpiryHours int
	}{
		ResetLink:   resetLink,
		ExpiryHours: 24,
	}

	return es.manager.NewBuilder().
		To(to).
		Subject("Password Reset Request").
		BodyFromTemplate("templates/password_reset.html", data).
		Tag("type", "password_reset").
		Tag("security", "true").
		Send(ctx)
}

// SendNotificationEmail sends a notification with fallback
func (es *EmailService) SendNotificationEmail(ctx context.Context, to, subject, message string) error {
	data := struct {
		Message string
		Subject string
	}{
		Message: message,
		Subject: subject,
	}

	email, err := es.manager.NewBuilder().
		To(to).
		Subject(subject).
		BodyFromTemplate("templates/notification.html", data).
		Tag("type", "notification").
		Build()

	if err != nil {
		return err
	}

	return es.manager.SendWithFallback(ctx, email)
}

// SendBulkEmail sends emails to multiple recipients using different providers for load distribution
func (es *EmailService) SendBulkEmail(ctx context.Context, recipients []string, subject, htmlBody string) []error {
	var errors []error
	providers := es.manager.GetAvailableProviders()

	for i, recipient := range recipients {
		// Distribute load across providers
		providerIndex := i % len(providers)
		provider := providers[providerIndex]

		builder, err := es.manager.NewBuilderWithProvider(provider)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to get provider %s for %s: %w", provider, recipient, err))
			continue
		}

		err = builder.
			To(recipient).
			Subject(subject).
			BodyHTML(htmlBody).
			Tag("type", "bulk").
			Tag("provider", string(provider)).
			Send(ctx)

		if err != nil {
			errors = append(errors, fmt.Errorf("failed to send to %s via %s: %w", recipient, provider, err))
		}
	}

	return errors
}

// GetProviderStats returns statistics about configured providers
func (es *EmailService) GetProviderStats() map[Provider]map[string]interface{} {
	stats := make(map[Provider]map[string]interface{})

	for _, provider := range es.manager.GetAvailableProviders() {
		stats[provider] = map[string]interface{}{
			"configured": true,
			"type":       string(provider),
		}
	}

	return stats
}

// SwitchPrimaryProvider switches the primary email provider
func (es *EmailService) SwitchPrimaryProvider(provider Provider) error {
	return es.manager.SetPrimaryProvider(provider)
}

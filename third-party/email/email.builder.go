package email

import (
	"context"
	"fmt"
)

func NewEmailBuilder(sender EmailSender) EmailBuilder {
	return &emailBuilder{
		email:  &Email{Tags: make(map[string]string)},
		sender: sender,
		errors: []error{},
	}
}

func (b *emailBuilder) To(to string) EmailBuilder {
	if to == "" {
		b.errors = append(b.errors, fmt.Errorf("recipient email cannot be empty"))
		return b
	}
	b.email.To = to
	return b
}

func (b *emailBuilder) Subject(subject string) EmailBuilder {
	if subject == "" {
		b.errors = append(b.errors, fmt.Errorf("subject cannot be empty"))
		return b
	}
	b.email.Subject = subject
	return b
}

func (b *emailBuilder) BodyHTML(body string) EmailBuilder {
	if body == "" {
		b.errors = append(b.errors, fmt.Errorf("HTML body cannot be empty"))
		return b
	}
	b.email.BodyHTML = body
	return b
}

func (b *emailBuilder) BodyText(body string) EmailBuilder {
	if body == "" {
		b.errors = append(b.errors, fmt.Errorf("text body cannot be empty"))
		return b
	}
	b.email.BodyText = body
	return b
}

func (b *emailBuilder) BodyFromTemplate(templatePath string, data interface{}) EmailBuilder {
	if templatePath == "" {
		b.errors = append(b.errors, fmt.Errorf("template path cannot be empty"))
		return b
	}
	if data == nil {
		b.errors = append(b.errors, fmt.Errorf("template data cannot be nil"))
		return b
	}
	b.email.TemplatePath = templatePath
	b.email.TemplateData = data
	return b
}

func (b *emailBuilder) CC(cc ...string) EmailBuilder {
	for _, addr := range cc {
		if addr == "" {
			b.errors = append(b.errors, fmt.Errorf("CC email cannot be empty"))
			return b
		}
	}
	b.email.CC = append(b.email.CC, cc...)
	return b
}

func (b *emailBuilder) BCC(bcc ...string) EmailBuilder {
	for _, addr := range bcc {
		if addr == "" {
			b.errors = append(b.errors, fmt.Errorf("BCC email cannot be empty"))
			return b
		}
	}
	b.email.BCC = append(b.email.BCC, bcc...)
	return b
}

func (b *emailBuilder) ReplyTo(replyTo string) EmailBuilder {
	if replyTo == "" {
		b.errors = append(b.errors, fmt.Errorf("reply-to email cannot be empty"))
		return b
	}
	b.email.ReplyTo = replyTo
	return b
}

func (b *emailBuilder) Tag(key, value string) EmailBuilder {
	if key == "" || value == "" {
		b.errors = append(b.errors, fmt.Errorf("tag key and value cannot be empty"))
		return b
	}
	b.email.Tags[key] = value
	return b
}

func (b *emailBuilder) Build() (*Email, error) {
	if len(b.errors) > 0 {
		return nil, fmt.Errorf("email validation failed: %v", b.errors)
	}

	if b.email.To == "" {
		return nil, fmt.Errorf("recipient email is required")
	}

	if b.email.Subject == "" {
		return nil, fmt.Errorf("subject is required")
	}

	if b.email.BodyHTML == "" && b.email.BodyText == "" && b.email.TemplatePath == "" {
		return nil, fmt.Errorf("email body, text, or template must be provided")
	}

	return b.email, nil
}

func (b *emailBuilder) Send(ctx context.Context) error {
	email, err := b.Build()
	if err != nil {
		return err
	}

	return b.sender.SendEmail(ctx, email)
}

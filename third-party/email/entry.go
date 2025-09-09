package email

func NewSMTPEmailBuilder(cfg Config) (EmailBuilder, error) {
	cfg.Provider = ProviderSMTP
	sender, err := NewEmailSender(cfg)
	if err != nil {
		return nil, err
	}
	return NewEmailBuilder(sender), nil
}

func NewSESEmailBuilder(cfg Config) (EmailBuilder, error) {
	cfg.Provider = ProviderSES
	sender, err := NewEmailSender(cfg)
	if err != nil {
		return nil, err
	}
	return NewEmailBuilder(sender), nil
}

func NewResendEmailBuilder(cfg Config) (EmailBuilder, error) {
	cfg.Provider = ProviderResend
	sender, err := NewEmailSender(cfg)
	if err != nil {
		return nil, err
	}
	return NewEmailBuilder(sender), nil
}

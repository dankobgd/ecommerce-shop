package config

// SetDefaults sets default values for AppSettings
func (s *AppSettings) SetDefaults() {
	if s.Host == "" {
		s.Host = "127.0.0.1"
	}
	if s.Port == 0 {
		s.Port = 3001
	}
	if s.ENV == "" {
		s.ENV = "development"
	}
}

// SetDefaults sets default values for DatabaseSettings
func (s *DatabaseSettings) SetDefaults() {
	if s.PostgresHost == "" {
		s.PostgresHost = "localhost"
	}
	if s.PostgresURI == "" {
		s.PostgresURI = ""
	}
	if s.PostgresDB == "" {
		s.PostgresDB = "ecommerce"
	}
	if s.PostgresPass == "" {
		s.PostgresPass = "test"
	}
	if s.PostgresUser == "" {
		s.PostgresUser = "test"
	}
}

// SetDefaults sets default values for AuthSettings
func (s *AuthSettings) SetDefaults() {
	if s.AccessTokenSecret == "" {
		s.AccessTokenSecret = "secret1"
	}
	if s.RefreshTokenSecret == "" {
		s.RefreshTokenSecret = "secret2"
	}
	if s.ResetPasswordValidFor == 0 {
		s.ResetPasswordValidFor = 0
	}
	if s.VerificationRequired == false {
		s.VerificationRequired = true
	}
}

// SetDefaults sets default values for EmailSettings
func (s *EmailSettings) SetDefaults() {
	if s.Enabled == false {
		s.Enabled = false
	}
	if s.From == "" {
		s.From = ""
	}
	if s.Host == "" {
		s.Host = ""
	}
	if s.Port == 0 {
		s.Port = 0
	}
	if s.User == "" {
		s.User = ""
	}
	if s.Pass == "" {
		s.Pass = ""
	}
	if s.Transport == "" {
		s.Transport = ""
	}
}

// SetDefaults sets default values for CookieSettings
func (s *CookieSettings) SetDefaults() {
	if s.Name == "" {
		s.Name = ""
	}
	if s.Path == "" {
		s.Path = ""
	}
	if s.Secret == "" {
		s.Secret = ""
	}
	if s.MaxAge == 0 {
		s.MaxAge = 0
	}
	if s.HTTPOnly == false {
		s.HTTPOnly = false
	}
	if s.Secure == false {
		s.Secure = false
	}
}

// SetDefaults sets default values for PasswordSettings
func (s *PasswordSettings) SetDefaults() {
	if s.MinLength == 0 {
		s.MinLength = 5
	}
	if s.MaxLength == 0 {
		s.MaxLength = 60
	}
	if s.Lowercase == false {
		s.Lowercase = true
	}
	if s.Uppercase == false {
		s.Uppercase = true
	}
	if s.Number == false {
		s.Number = true
	}
	if s.Symbol == false {
		s.Symbol = true
	}
}

// SetDefaults sets default values for LoggerSettings
func (s *LoggerSettings) SetDefaults() {
	if s.EnableConsole == false {
		s.EnableConsole = true
	}
	if s.ConsoleJSON == false {
		s.EnableConsole = false
	}
	if s.ConsoleLevel == "" {
		s.ConsoleLevel = "debug"
	}
	if s.EnableFile == false {
		s.EnableFile = false
	}
	if s.FileJSON == false {
		s.FileJSON = true
	}
	if s.FileLevel == "" {
		s.FileLevel = "info"
	}
	if s.FileLocation == "" {
		s.FileLocation = ""
	}
}

package language

// LanguageConfig represents the language configuration domain model
type LanguageConfig struct {
	Type   string                 `yaml:"type" json:"type"`
	Config map[string]interface{} `yaml:"config" json:"config"`
}

// SupportedLanguages defines the list of supported languages
var SupportedLanguages = []string{"go", "python", "nodejs", "java", "rust"}

// IsSupported checks if the language is supported
func (c *LanguageConfig) IsSupported() bool {
	for _, lang := range SupportedLanguages {
		if c.Type == lang {
			return true
		}
	}
	return false
}

// GetConfig retrieves a configuration value
func (c *LanguageConfig) GetConfig(key string) (interface{}, bool) {
	if c.Config == nil {
		return nil, false
	}
	val, ok := c.Config[key]
	return val, ok
}

// Validate validates the language configuration
func (c *LanguageConfig) Validate() error {
	if c.Type == "" {
		return ErrLanguageTypeRequired
	}

	if !c.IsSupported() {
		return NewErrUnsupportedLanguage(c.Type)
	}

	return nil
}

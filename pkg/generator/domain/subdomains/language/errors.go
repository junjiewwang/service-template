package language

import "fmt"

var (
	ErrLanguageTypeRequired = fmt.Errorf("language type is required")
)

type ErrUnsupportedLanguage struct {
	Language string
}

func (e *ErrUnsupportedLanguage) Error() string {
	return fmt.Sprintf("unsupported language: %s (supported: %v)", e.Language, SupportedLanguages)
}

func NewErrUnsupportedLanguage(lang string) error {
	return &ErrUnsupportedLanguage{Language: lang}
}

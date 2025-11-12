package events

// Generator event names
const (
	EventGenerationStarted   = "generation.started"
	EventGenerationCompleted = "generation.completed"
	EventGenerationFailed    = "generation.failed"
	EventValidationStarted   = "validation.started"
	EventValidationCompleted = "validation.completed"
	EventValidationFailed    = "validation.failed"
	EventTemplateRendered    = "template.rendered"
)

// GenerationStartedEvent is published when generation starts
type GenerationStartedEvent struct {
	BaseEvent
}

// NewGenerationStartedEvent creates a new generation started event
func NewGenerationStartedEvent(generatorName string, data map[string]interface{}) *GenerationStartedEvent {
	return &GenerationStartedEvent{
		BaseEvent: NewBaseEvent(EventGenerationStarted, generatorName, data),
	}
}

// GenerationCompletedEvent is published when generation completes successfully
type GenerationCompletedEvent struct {
	BaseEvent
}

// NewGenerationCompletedEvent creates a new generation completed event
func NewGenerationCompletedEvent(generatorName string, data map[string]interface{}) *GenerationCompletedEvent {
	return &GenerationCompletedEvent{
		BaseEvent: NewBaseEvent(EventGenerationCompleted, generatorName, data),
	}
}

// GenerationFailedEvent is published when generation fails
type GenerationFailedEvent struct {
	BaseEvent
}

// NewGenerationFailedEvent creates a new generation failed event
func NewGenerationFailedEvent(generatorName string, data map[string]interface{}) *GenerationFailedEvent {
	return &GenerationFailedEvent{
		BaseEvent: NewBaseEvent(EventGenerationFailed, generatorName, data),
	}
}

// ValidationStartedEvent is published when validation starts
type ValidationStartedEvent struct {
	BaseEvent
}

// NewValidationStartedEvent creates a new validation started event
func NewValidationStartedEvent(generatorName string, data map[string]interface{}) *ValidationStartedEvent {
	return &ValidationStartedEvent{
		BaseEvent: NewBaseEvent(EventValidationStarted, generatorName, data),
	}
}

// ValidationCompletedEvent is published when validation completes
type ValidationCompletedEvent struct {
	BaseEvent
}

// NewValidationCompletedEvent creates a new validation completed event
func NewValidationCompletedEvent(generatorName string, data map[string]interface{}) *ValidationCompletedEvent {
	return &ValidationCompletedEvent{
		BaseEvent: NewBaseEvent(EventValidationCompleted, generatorName, data),
	}
}

// ValidationFailedEvent is published when validation fails
type ValidationFailedEvent struct {
	BaseEvent
}

// NewValidationFailedEvent creates a new validation failed event
func NewValidationFailedEvent(generatorName string, data map[string]interface{}) *ValidationFailedEvent {
	return &ValidationFailedEvent{
		BaseEvent: NewBaseEvent(EventValidationFailed, generatorName, data),
	}
}

// TemplateRenderedEvent is published when a template is rendered
type TemplateRenderedEvent struct {
	BaseEvent
}

// NewTemplateRenderedEvent creates a new template rendered event
func NewTemplateRenderedEvent(generatorName string, data map[string]interface{}) *TemplateRenderedEvent {
	return &TemplateRenderedEvent{
		BaseEvent: NewBaseEvent(EventTemplateRendered, generatorName, data),
	}
}

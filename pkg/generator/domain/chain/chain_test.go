package chain

import (
	"context"
	"testing"
)

// Mock handlers for testing
type mockParserHandler struct {
	*BaseHandler
	executed bool
}

func newMockParserHandler(name string) *mockParserHandler {
	return &mockParserHandler{
		BaseHandler: NewBaseHandler(name),
	}
}

func (h *mockParserHandler) Handle(ctx *ProcessingContext) error {
	h.executed = true
	return h.CallNext(ctx)
}

func (h *mockParserHandler) Parse(ctx *ProcessingContext) error {
	return h.Handle(ctx)
}

type mockValidatorHandler struct {
	*BaseHandler
	executed bool
}

func newMockValidatorHandler(name string) *mockValidatorHandler {
	return &mockValidatorHandler{
		BaseHandler: NewBaseHandler(name),
	}
}

func (h *mockValidatorHandler) Handle(ctx *ProcessingContext) error {
	h.executed = true
	return h.CallNext(ctx)
}

func (h *mockValidatorHandler) Validate(ctx *ProcessingContext) error {
	return h.Handle(ctx)
}

type mockGeneratorHandler struct {
	*BaseHandler
	executed bool
}

func newMockGeneratorHandler(name string) *mockGeneratorHandler {
	return &mockGeneratorHandler{
		BaseHandler: NewBaseHandler(name),
	}
}

func (h *mockGeneratorHandler) Handle(ctx *ProcessingContext) error {
	h.executed = true
	return h.CallNext(ctx)
}

func (h *mockGeneratorHandler) Generate(ctx *ProcessingContext) error {
	return h.Handle(ctx)
}

// TestChainBuilder tests the chain builder
func TestChainBuilder(t *testing.T) {
	t.Run("Build empty chain", func(t *testing.T) {
		builder := NewChainBuilder()
		chain := builder.Build()

		if chain != nil {
			t.Error("Expected nil for empty chain")
		}
	})

	t.Run("Build chain with handlers", func(t *testing.T) {
		h1 := newMockParserHandler("handler1")
		h2 := newMockParserHandler("handler2")
		h3 := newMockParserHandler("handler3")

		builder := NewChainBuilder()
		chain := builder.Add(h1).Add(h2).Add(h3).Build()

		if chain == nil {
			t.Fatal("Expected non-nil chain")
		}

		ctx := NewProcessingContext(context.Background(), nil)
		err := chain.Handle(ctx)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if !h1.executed || !h2.executed || !h3.executed {
			t.Error("Not all handlers were executed")
		}
	})
}

// TestChainExecutor tests the chain executor
func TestChainExecutor(t *testing.T) {
	t.Run("Execute chain", func(t *testing.T) {
		h1 := newMockParserHandler("handler1")
		h2 := newMockParserHandler("handler2")

		builder := NewChainBuilder()
		chain := builder.Add(h1).Add(h2).Build()

		executor := NewChainExecutor(chain)
		ctx := NewProcessingContext(context.Background(), nil)

		err := executor.Execute(ctx)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if !h1.executed || !h2.executed {
			t.Error("Not all handlers were executed")
		}
	})

	t.Run("Execute nil chain", func(t *testing.T) {
		executor := NewChainExecutor(nil)
		ctx := NewProcessingContext(context.Background(), nil)

		err := executor.Execute(ctx)
		if err == nil {
			t.Error("Expected error for nil chain")
		}
	})
}

// TestProcessingContext tests the processing context
func TestProcessingContext(t *testing.T) {
	t.Run("Domain model operations", func(t *testing.T) {
		ctx := NewProcessingContext(context.Background(), nil)

		// Set and get domain model
		ctx.SetDomainModel("service", map[string]string{"name": "test"})

		model, ok := ctx.GetDomainModel("service")
		if !ok {
			t.Error("Expected to find domain model")
		}

		modelMap := model.(map[string]string)
		if modelMap["name"] != "test" {
			t.Error("Domain model value mismatch")
		}
	})

	t.Run("Error operations", func(t *testing.T) {
		ctx := NewProcessingContext(context.Background(), nil)

		if ctx.HasErrors() {
			t.Error("Expected no errors initially")
		}

		ctx.AddError(context.Canceled)

		if !ctx.HasErrors() {
			t.Error("Expected to have errors")
		}

		errors := ctx.GetErrors()
		if len(errors) != 1 {
			t.Errorf("Expected 1 error, got %d", len(errors))
		}
	})

	t.Run("Validation error operations", func(t *testing.T) {
		ctx := NewProcessingContext(context.Background(), nil)

		ctx.AddValidationError("service", context.Canceled)
		ctx.AddValidationError("service", context.DeadlineExceeded)

		errors := ctx.GetValidationErrors("service")
		if len(errors) != 2 {
			t.Errorf("Expected 2 validation errors, got %d", len(errors))
		}
	})

	t.Run("Generated file operations", func(t *testing.T) {
		ctx := NewProcessingContext(context.Background(), nil)

		content := []byte("test content")
		ctx.AddGeneratedFile("/path/to/file", content)

		retrieved, ok := ctx.GetGeneratedFile("/path/to/file")
		if !ok {
			t.Error("Expected to find generated file")
		}

		if string(retrieved) != string(content) {
			t.Error("Generated file content mismatch")
		}
	})

	t.Run("Metadata operations", func(t *testing.T) {
		ctx := NewProcessingContext(context.Background(), nil)

		ctx.SetMetadata("key", "value")

		value, ok := ctx.GetMetadata("key")
		if !ok {
			t.Error("Expected to find metadata")
		}

		if value.(string) != "value" {
			t.Error("Metadata value mismatch")
		}
	})
}

// TestMiddleware tests middleware functionality
func TestMiddleware(t *testing.T) {
	t.Run("Logging middleware", func(t *testing.T) {
		handler := newMockParserHandler("test")
		wrapped := NewLoggingMiddleware(handler)

		ctx := NewProcessingContext(context.Background(), nil)
		err := wrapped.Handle(ctx)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if !handler.executed {
			t.Error("Handler was not executed")
		}
	})

	t.Run("Metrics middleware", func(t *testing.T) {
		handler := newMockParserHandler("test")
		wrapped := NewMetricsMiddleware(handler)

		ctx := NewProcessingContext(context.Background(), nil)
		err := wrapped.Handle(ctx)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if !handler.executed {
			t.Error("Handler was not executed")
		}

		metricsHandler := wrapped.(*MetricsMiddleware)
		metrics := metricsHandler.GetMetrics()

		if metrics.ExecutionCount != 1 {
			t.Errorf("Expected execution count 1, got %d", metrics.ExecutionCount)
		}
	})

	t.Run("Recovery middleware", func(t *testing.T) {
		handler := newMockParserHandler("test")
		wrapped := NewRecoveryMiddleware(handler)

		ctx := NewProcessingContext(context.Background(), nil)
		err := wrapped.Handle(ctx)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}

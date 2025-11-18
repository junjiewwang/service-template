package chain

import "fmt"

// ChainExecutor executes a chain of handlers
type ChainExecutor struct {
	chain Handler
}

// NewChainExecutor creates a new chain executor
func NewChainExecutor(chain Handler) *ChainExecutor {
	return &ChainExecutor{
		chain: chain,
	}
}

// Execute executes the chain with the given context
func (e *ChainExecutor) Execute(ctx *ProcessingContext) error {
	if e.chain == nil {
		return fmt.Errorf("chain is not initialized")
	}

	return e.chain.Handle(ctx)
}

// ExecuteWithRecovery executes the chain with panic recovery
func (e *ChainExecutor) ExecuteWithRecovery(ctx *ProcessingContext) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered: %v", r)
		}
	}()

	return e.Execute(ctx)
}

package resolvers

import "context"

// KeepExistingResolverID is the unique identifier for the keep existing resolver
const KeepExistingResolverID = "keep_existing"

// KeepExistingResolver keeps the existing content when conflict occurs
type KeepExistingResolver struct{}

func init() {
	DefaultResolverRegistry.MustRegister(&KeepExistingResolver{})
}

// ID returns the resolver identifier
func (r *KeepExistingResolver) ID() string {
	return KeepExistingResolverID
}

// Description returns the resolver description
func (r *KeepExistingResolver) Description() string {
	return "Keep existing content when conflict occurs (preserve user modifications)"
}

// Resolve returns the existing content, preserving user modifications
func (r *KeepExistingResolver) Resolve(ctx context.Context, input *ConflictInput) ([]byte, error) {
	return input.ExistingContent, nil
}

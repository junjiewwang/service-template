package resolvers

import "context"

// UseNewResolverID is the unique identifier for the use new resolver
const UseNewResolverID = "use_new"

// UseNewResolver uses the new content when conflict occurs
type UseNewResolver struct{}

func init() {
	DefaultResolverRegistry.MustRegister(&UseNewResolver{})
}

// ID returns the resolver identifier
func (r *UseNewResolver) ID() string {
	return UseNewResolverID
}

// Description returns the resolver description
func (r *UseNewResolver) Description() string {
	return "Use new content when conflict occurs (overwrite user modifications)"
}

// Resolve returns the new content, overwriting user modifications
func (r *UseNewResolver) Resolve(ctx context.Context, input *ConflictInput) ([]byte, error) {
	return input.NewContent, nil
}

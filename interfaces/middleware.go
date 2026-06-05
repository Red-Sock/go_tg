package interfaces

// Middleware wraps a Handler. Call next.Handle to continue the chain.
// Return without calling next to abort (e.g. access denied, feature flag off).
type Middleware func(next Handler) Handler

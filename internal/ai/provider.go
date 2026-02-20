package ai

import "context"

// Provider abstracts an LLM backend. Implementations must support stateful
// conversation threading — each call may reference prior context.
type Provider interface {
	// Send sends a user message and returns the assistant reply.
	// systemPrompt is injected as instructions when non-empty.
	// The provider is responsible for maintaining conversational state.
	Send(ctx context.Context, message, systemPrompt string) (string, error)
}

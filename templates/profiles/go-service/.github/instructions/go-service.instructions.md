---
name: Go Service
description: Composable Go services with clear interfaces, explicit errors, and operational hygiene
applyTo: "**/*.go"
---

# Go

Go is great when you need it — when you need performance, small binaries,
excellent concurrency, and operational simplicity. These rules help you write
Go that stays simple as codebases grow.

## Structure

Organize by domain, not by pattern:

```
cmd/
  server/
    main.go              # Entry point — wiring only
internal/
  order/
    order.go             # Domain types and business logic
    service.go           # Application service (orchestration)
    postgres.go          # Repository implementation
    handler.go           # HTTP handler
  customer/
    customer.go
    service.go
    handler.go
pkg/
  middleware/            # Shared HTTP middleware
  nats/                  # NATS client helpers
```

- **`cmd/` is thin.** It wires dependencies and starts the server. No business logic.
- **`internal/` is your domain.** Organized by feature, not by layer.
- **`pkg/` is truly reusable.** Only put code here that other projects could import.

## Interfaces

Interfaces in Go are powerful *because* they're small. Keep them that way.

```go
// ✅ Small, focused interface — defined by the consumer
type OrderRepository interface {
    FindByID(ctx context.Context, id string) (*Order, error)
    Save(ctx context.Context, order *Order) error
}

// ❌ Kitchen sink interface
type OrderRepository interface {
    FindByID(ctx context.Context, id string) (*Order, error)
    FindAll(ctx context.Context) ([]*Order, error)
    FindByCustomer(ctx context.Context, customerID string) ([]*Order, error)
    Save(ctx context.Context, order *Order) error
    Delete(ctx context.Context, id string) error
    Count(ctx context.Context) (int, error)
    // ... 12 more methods
}
```

- **Define interfaces where they're used**, not where they're implemented.
  The consumer knows what it needs; the implementation doesn't need to know
  who's calling.
- **Accept interfaces, return structs.** Callers should depend on behavior;
  implementations should return concrete types.

## Error handling

Errors are values. Treat them with care.

```go
// ✅ Errors with context — helpful when debugging at 3am
if err != nil {
    return fmt.Errorf("loading order %s: %w", orderID, err)
}

// ❌ Swallowed error — silent corruption
if err != nil {
    log.Println(err)
}
```

- Always wrap errors with context using `fmt.Errorf("what happened: %w", err)`.
- Use sentinel errors (`var ErrNotFound = errors.New("not found")`) for
  conditions callers need to check.
- Use custom error types when errors need to carry structured data.
- **Never ignore errors.** If you truly don't care, document why with a comment.

## Concurrency

- **Always pass `context.Context` as the first parameter.** Respect cancellation.
- Use `errgroup.Group` for coordinating concurrent work with error collection.
- Protect shared state with `sync.Mutex` when needed, but prefer channels for
  communication between goroutines.
- **Don't start goroutines you can't stop.** Every goroutine should have a clear
  shutdown path.

```go
// ✅ Controlled concurrency with errgroup
g, ctx := errgroup.WithContext(ctx)
for _, item := range items {
    g.Go(func() error {
        return processItem(ctx, item)
    })
}
if err := g.Wait(); err != nil {
    return fmt.Errorf("processing items: %w", err)
}
```

## Testing

- Use table-driven tests for functions with many cases.
- Use `testify/assert` or `testify/require` for readable assertions.
- Test behavior at the handler level with `httptest.NewRequest`.
- Use `t.Parallel()` for independent tests.
- Mock interfaces, not functions. Keep mocks close to tests.

## What to avoid

- Premature abstractions — write concrete code first, extract when patterns emerge.
- Package-level `init()` functions — explicit initialization is clearer.
- `interface{}` / `any` when you can use generics.
- Global variables for configuration — pass config explicitly.
- Frameworks that hide Go's simplicity. The stdlib is excellent — use it.
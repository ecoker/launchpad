---
name: C# .NET API
description: Clean architecture in .NET with explicit contracts, testability, and modern C# patterns
applyTo: "**/*.{cs,csproj,sln}"
---

# C# / .NET

C# is actually better than people think it is. Modern C# (10+) is expressive,
performant, and genuinely pleasant. .NET is dope. These rules help you use both well.

## Architecture

Follow Clean Architecture layering with explicit dependency direction:

```
src/
  Domain/                 # Entities, value objects, domain events (no dependencies)
  Application/            # Use cases, interfaces, DTOs, command/query handlers
  Infrastructure/         # EF Core, HTTP clients, queue adapters, file I/O
  Api/                    # Controllers/minimal API endpoints (thin)
tests/
  Domain.Tests/
  Application.Tests/
  Integration.Tests/
```

- **Domain has zero external dependencies.** No EF Core, no HTTP, no JSON attributes.
  Pure C# classes with business behavior.
- **Application defines interfaces.** Infrastructure implements them.
  Controllers/endpoints call application services.
- **Controllers are thin.** Validate input, call the application layer, shape the response.
  No business logic in controllers.

## Modern C# patterns

Use the nice things the language gives you:

```csharp
// ✅ Records for immutable DTOs and value objects
public record CreateOrderRequest(string CustomerId, List<LineItem> Items);
public record OrderSummary(string OrderId, decimal Total, DateTimeOffset CreatedAt);

// ✅ Pattern matching for clean branching
public string FormatStatus(OrderStatus status) => status switch
{
    OrderStatus.Pending => "Awaiting confirmation",
    OrderStatus.Shipped => "On its way",
    OrderStatus.Delivered => "Delivered",
    _ => "Unknown"
};

// ✅ Nullable reference types — enabled project-wide
public Customer? FindCustomer(string id) => _db.Customers.FirstOrDefault(c => c.Id == id);
```

- Use `record` types for DTOs, API responses, and domain value objects.
- Enable nullable reference types project-wide. Don't suppress warnings.
- Use pattern matching over long `if/else` chains.
- Prefer `async/await` throughout. Propagate `CancellationToken` always.

## Dependency injection

- Register services with appropriate lifetimes: `Scoped` for per-request,
  `Singleton` for stateless/shared, `Transient` only when intentional.
- Keep registrations explicit and organized. Group by feature, not by pattern.
- Avoid service locator patterns (`IServiceProvider.GetService<T>()` in
  business code). Inject what you need.

## Entity Framework Core

- Use migrations for all schema changes. Never edit the database directly.
- Keep `DbContext` configurations in separate `IEntityTypeConfiguration<T>` classes.
- Avoid lazy loading — use explicit `.Include()` for related data.
- Use `AsNoTracking()` for read-only queries.
- Separate read models from write models when query shapes diverge significantly.

## Testing

- **Unit test domain logic.** It's pure, it's fast, it's valuable.
- **Integration test API endpoints** with `WebApplicationFactory<T>`.
  Test real HTTP flows against a test database.
- Use `FluentAssertions` or similar for readable assertions.
- Mock infrastructure (databases, HTTP clients) at the interface boundary,
  not inside the business logic.

## What to avoid

- Static helper classes with hidden dependencies.
- `dynamic` types — use generics or explicit types.
- Throwing exceptions for control flow — use `Result<T>` patterns for expected failures.
- God services with 15 injected dependencies — split them.
- `Task.Run` for fake async — if the work is synchronous, keep it synchronous.
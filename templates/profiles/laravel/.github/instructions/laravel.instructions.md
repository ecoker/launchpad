---
name: Laravel
description: Laravel conventions with domain clarity, service boundaries, and testing discipline
applyTo: "**/*.{php,blade.php}"
---

# Laravel

If you have to do PHP, use Laravel. And if you use Laravel, do it well.
Laravel's conventions are powerful — lean into them, but don't let convenience
become chaos.

## Architecture

Laravel gives you an opinionated structure. Extend it, don't fight it:

```
app/
  Models/               # Eloquent models — focused, validated
  Http/
    Controllers/        # Thin. Validate, delegate, respond.
    Requests/           # Form request validation — the boundary layer
    Resources/          # API resource transformations
  Services/             # Business logic that doesn't belong in models
  Actions/              # Single-purpose classes for complex operations
  Events/               # Domain events
  Listeners/            # Event handlers (side effects)
  Policies/             # Authorization logic
database/
  migrations/           # Versioned, reviewable, reversible
tests/
  Feature/              # HTTP + integration tests
  Unit/                 # Pure logic tests
```

## Controllers

Controllers are routers, not business logic containers.

```php
// ✅ Thin controller — delegates to an action
class CreateOrderController extends Controller
{
    public function __invoke(CreateOrderRequest $request, CreateOrder $action)
    {
        $order = $action->execute($request->validated());
        return new OrderResource($order);
    }
}

// ❌ Fat controller — validation, logic, queries, emails all in one
class OrderController extends Controller
{
    public function store(Request $request)
    {
        // 80 lines of validation, saving, emailing, logging...
    }
}
```

- Use **invokable controllers** for single-action endpoints.
- Always use **Form Requests** for validation — not inline `$request->validate()`.
- Return **API Resources** for consistent JSON shaping.

## Eloquent

- Keep models focused on relationships, scopes, and accessors.
- **Don't put business logic in models.** Use Services or Actions.
- Use **scopes** for reusable query conditions.
- Use **casts** and **accessors** for type safety on attributes.

```php
// ✅ Clean model with scopes and casts
class Order extends Model
{
    protected $casts = [
        'total_cents' => 'integer',
        'shipped_at' => 'datetime',
    ];

    public function scopeActive(Builder $query): Builder
    {
        return $query->whereNull('cancelled_at');
    }

    public function customer(): BelongsTo
    {
        return $this->belongsTo(Customer::class);
    }
}
```

## Services and Actions

- **Services** orchestrate multiple steps with dependencies.
- **Actions** are single-purpose, invokable classes — great for complex
  operations that need to be reusable across controllers, commands, and jobs.
- Both should be explicitly typed with clear inputs and outputs.

## Testing

- Write **Feature tests** for every endpoint. Test the full HTTP cycle.
- Use **factories and seeders** for test data — never hardcode IDs.
- Use `RefreshDatabase` trait for clean test isolation.
- Test authorization (policies) explicitly — don't assume middleware coverage.

## What to avoid

- Raw DB queries when Eloquent works fine.
- `dd()` and `dump()` in committed code.
- Forgetting to authorize — every controller action checks authorization.
- Massive service providers with manual bindings — use auto-discovery.
- Storing business logic in Blade templates.
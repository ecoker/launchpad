---
name: Server-Side Patterns
description: Validation, error handling, data access, and form conventions for backend code
applyTo: "**/*.{ts,tsx,js,jsx,ex,exs,py,go,rs,rb,cs,java,kt,php}"
---

# Server-side patterns

Every backend surface — API endpoints, form actions, background jobs, message
handlers — follows the same discipline: validate at the boundary, handle errors
explicitly, and keep data access predictable.

## Validation

Validate **at the boundary**, not deep inside domain logic. Every piece of
external input (HTTP request, message payload, file upload, environment
variable) must be validated and typed before it touches your domain.

### The rule

```
External input → Validate + Parse → Typed domain object → Business logic
```

Never pass raw request bodies, untyped maps, or `any`/`dynamic` values into
domain functions. Parse, don't validate — turn untrusted input into trusted
types at the edge.

### Framework-specific validation

| Stack | Validation approach |
|-------|-------------------|
| **TypeScript (Next.js)** | Zod schemas in Server Actions and API routes. Parse with `.safeParse()`, return structured errors. |
| **TypeScript (SvelteKit)** | Zod or Valibot in `+page.server.ts` actions and load functions. |
| **TypeScript (Fastify)** | TypeBox schemas on route definitions — runtime validation + TypeScript inference from one source. |
| **TypeScript (React Router v7)** | Zod in action/loader functions. Validate before touching domain logic. |
| **Elixir (Phoenix)** | Ecto changesets with explicit validations. Embedded schemas for non-persisted input. |
| **Python (FastAPI)** | Pydantic models as function parameters — automatic validation + serialization + OpenAPI docs. |
| **Python (Django)** | Django Forms or DRF Serializers. `form.is_valid()` before any processing. |
| **Go** | Struct tags with `go-playground/validator` or manual validation functions. Return typed errors. |
| **Rust (Axum)** | `serde` for deserialization + `validator` crate or manual checks. Return typed error responses. |
| **Ruby (Rails)** | Strong Parameters in controllers. Model validations for persistence. Form objects for complex input. |
| **C# (.NET)** | Data Annotations on DTOs + `[ApiController]` auto-validation. FluentValidation for complex rules. |
| **Java (Spring)** | `@Valid` + Bean Validation annotations. Custom `ConstraintValidator` for domain rules. |
| **Laravel** | Form Request classes with `rules()` method. Always validate — never trust `$request->input()` directly. |
| **Dart (Flutter)** | Validate at the network boundary. Use `freezed`/`json_serializable` for typed deserialization. Client-side form validation via `Form` + `TextFormField.validator`. |

### Validation errors shape

Return validation errors in a consistent, machine-readable format. Every
endpoint should return errors the frontend can map to specific fields:

```json
{
  "errors": {
    "email": ["must be a valid email address"],
    "age": ["must be at least 18"]
  }
}
```

For non-field errors (authorization, rate limits, system errors), use a
top-level `message` field:

```json
{
  "message": "You don't have permission to perform this action"
}
```

## Error handling

Errors are not exceptions to the happy path — they're expected outcomes that
deserve explicit modeling.

### Principles

- **Use the language's error modeling idiom.** Don't fight the language.
  - TypeScript: discriminated unions (`{ ok: true, data } | { ok: false, error }`)
  - Elixir: `{:ok, value} | {:error, reason}` tuples
  - Go: `(value, error)` return pairs
  - Rust: `Result<T, E>`
  - Python: raise specific exceptions, catch at the boundary
  - C#/Java: throw specific exceptions, catch at the controller/filter level
  - Ruby: raise specific errors, rescue at the controller level
- **Never swallow errors silently.** If you catch an error, either handle it
  meaningfully or re-raise it with context.
- **Categorize errors for the caller.** The client needs to know: was this
  their fault (4xx)? Our fault (5xx)? Transient (retry)? Permanent (stop)?

### HTTP error responses

Use consistent HTTP status codes and response shapes:

| Situation | Status | Response body |
|-----------|--------|--------------|
| Invalid input | `400` or `422` | Field-level errors object |
| Not authenticated | `401` | `{ "message": "..." }` |
| Not authorized | `403` | `{ "message": "..." }` |
| Resource not found | `404` | `{ "message": "..." }` |
| Conflict / duplicate | `409` | `{ "message": "..." }` |
| Server error | `500` | `{ "message": "..." }` (no stack traces in production) |

### Domain errors

Model domain errors as explicit types, not strings:

```typescript
// ✅ Explicit error type
type OrderError =
  | { kind: "out_of_stock"; productId: string }
  | { kind: "payment_declined"; reason: string }
  | { kind: "already_shipped" };

// ❌ Stringly-typed error
throw new Error("out of stock");
```

```elixir
# ✅ Tagged tuples
{:error, :out_of_stock}
{:error, {:payment_declined, reason}}

# ❌ Bare strings
{:error, "out of stock"}
```

### Centralized error handling

Use your framework's error handling middleware to catch unhandled errors and
produce consistent responses:

- **Next.js**: `error.tsx` boundary components and try/catch in Server Actions
- **SvelteKit**: `+error.svelte` pages and `handleError` hook
- **Phoenix**: `ErrorView` and `FallbackController`
- **FastAPI**: `@app.exception_handler` decorators
- **Django**: custom middleware and `handler404`/`handler500`
- **Rails**: `rescue_from` in `ApplicationController`
- **Go**: middleware that recovers panics and formats error responses
- **Rust/Axum**: `IntoResponse` impl on custom error enum
- **.NET**: `IExceptionHandler` or `UseExceptionHandler` middleware
- **Spring**: `@RestControllerAdvice` with `@ExceptionHandler` methods
- **Laravel**: `Handler::render()` in `Exceptions/Handler.php`

## Form actions and mutations

For full-stack frameworks with server-rendered UI, mutations follow the
**action pattern**: form submission → server validation → redirect or
re-render with errors.

### Framework conventions

| Stack | Mutation pattern |
|-------|-----------------|
| **Next.js** | Server Actions (`"use server"`) — validate with Zod, return `{ errors }` or redirect |
| **SvelteKit** | Form actions in `+page.server.ts` — `fail(422, { errors })` or `redirect(303, url)` |
| **Phoenix LiveView** | `handle_event` with changeset validation → assign errors to socket |
| **Rails** | Controller actions with strong params → model validation → render or redirect |
| **Django** | View functions/classes with Form validation → render or redirect |
| **Laravel** | Controller methods with Form Request validation → redirect with errors |

### The pattern

1. Receive form data / request body
2. Validate immediately — before any business logic
3. On validation failure: return/render errors mapped to fields
4. On success: perform the action, then redirect (POST-Redirect-GET)
5. On domain error: return a meaningful error the UI can display

## Data access

- **Repositories / query modules at the boundary.** Domain logic calls a data
  access interface; it never constructs SQL or ORM queries directly.
- **Avoid N+1 queries.** Use eager loading (`includes`, `preload`,
  `select_related`, `Include()`, `JOIN FETCH`) deliberately.
- **Transactions for multi-step mutations.** If two writes must succeed
  together, wrap them in a transaction (`Ecto.Multi`, `@Transactional`,
  `ActiveRecord::Base.transaction`, `prisma.$transaction`).
- **Read-only queries should be explicit.** Mark read paths as read-only
  (`AsNoTracking()`, `Repo.all()` without preloading unnecessary assocs,
  `readonly` query hints) to avoid accidental writes and improve performance.

## Background jobs

- **Idempotent handlers.** Jobs will be retried. Design them so processing the
  same job twice produces the same result.
- **Typed payloads.** Job arguments should be validated/parsed at the boundary,
  just like HTTP input.
- **Timeouts and dead-letter queues.** Every job should have a timeout. Failed
  jobs should go somewhere observable, not disappear.

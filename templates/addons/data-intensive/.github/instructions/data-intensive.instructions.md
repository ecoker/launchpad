---
name: Data-Intensive Systems
description: Patterns for building reliable, observable, data-heavy applications
applyTo: "**/*.{sql,py,ts,tsx,js,jsx,go,cs,ex,exs,yaml,yml,json}"
---

# Data-intensive systems

> "Data outlives code." — Martin Kleppmann

If your system processes, stores, or moves meaningful amounts of data, these
rules apply. They're drawn heavily from *Designing Data-Intensive Applications*
and from hard-won production experience.

## Schema and contracts

- **Schema is a contract.** Treat it with the same rigor as a public API.
  Version it, review it, and never break it silently.
- **Migrations are code.** They live in version control, they're reviewed in PRs,
  and they're reversible (or at minimum, forward-compatible).
- **Use explicit types.** In every language: Pydantic models, TypeScript types,
  Go structs, Ecto schemas. Never pass raw `dict` / `any` / `map[string]interface{}`
  across boundaries.

```sql
-- ✅ Explicit, versioned, reviewable
ALTER TABLE orders ADD COLUMN shipped_at timestamptz;

-- ❌ Hand-editing production at 2am
-- UPDATE orders SET shipped = true WHERE id = 42;
```

## PostgreSQL

Postgres is our primary relational store. It's battle-tested, featureful, and
correct.

- Use `timestamptz` for all timestamps — never `timestamp` without timezone.
- Use `text` over `varchar(n)` unless you have a real constraint reason.
- Use `jsonb` columns sparingly and intentionally — if you're querying into JSON
  regularly, it probably wants to be a proper column or table.
- Index deliberately. Explain your queries. Avoid sequential scans on large tables.
- Use advisory locks or `SELECT ... FOR UPDATE` for critical sections — never
  application-level locking against a shared database.

## Message systems and NATS

- **NATS** is our go-to for lightweight pub-sub, request-reply, and work queues.
  It's fast, operationally simple, and has excellent multi-language support.
- Design message schemas explicitly. A message is a contract between services.
- Consumers must be **idempotent**. Messages will be delivered at-least-once.
  Design handlers so that processing the same message twice produces the same result.
- Use JetStream for durable, exactly-once delivery semantics when needed.

```
// ✅ Idempotent handler — safe to retry
async function handleOrderShipped(msg: OrderShippedEvent) {
  await db.query(
    `UPDATE orders SET shipped_at = $1 WHERE id = $2 AND shipped_at IS NULL`,
    [msg.shippedAt, msg.orderId]
  );
}

// ❌ Non-idempotent — double-processing causes duplicate charges
async function handleOrderShipped(msg: OrderShippedEvent) {
  await chargeCustomer(msg.orderId);
}
```

## Data formats

- **Parquet** is the default for analytical and archival data. It's columnar,
  compressed, and universally supported by analytics tooling.
- Use **JSON lines (JSONL)** for streaming/log data when human readability matters.
- Use **CSV** only for interchange with non-technical consumers. Always include headers.

## Observability

This is not optional. Every data system must be observable from day one.

- **Structured logs.** JSON, with consistent fields: `level`, `message`, `correlation_id`,
  `service`, `duration_ms`. Never `console.log("something happened")`.
- **Correlation IDs.** Every request and message gets an ID that flows through the
  entire processing chain. When something fails, you can trace it end-to-end.
- **Health checks.** Every service exposes a health endpoint that verifies its
  dependencies (database, queue, upstream services).
- **Metrics.** Track latency, throughput, error rates, and queue depth. If you
  can't measure it, you can't improve it.

## Pipeline design

- Break data transformations into **small, composable, testable steps**.
- Each step should be independently runnable and verifiable.
- Design for **replay safety**: you should be able to re-run any pipeline step
  from its input without side effects or data loss.
- Separate **extraction**, **transformation**, and **loading** concerns clearly.
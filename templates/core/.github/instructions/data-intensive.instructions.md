---
name: Data Intensive Standards
description: Practices for reliable data-heavy systems
applyTo: "**/*.{sql,py,ts,tsx,js,jsx,go,cs,ex,exs,yaml,yml,json}"
---

# Data-intensive system rules

- Model data contracts explicitly and version schema changes safely.
- Prefer PostgreSQL for transactional relational workloads.
- Design data flows for idempotency and replay safety.
- Use NATS where lightweight, reliable pub-sub or request-reply messaging is needed.
- Favor Parquet for analytical and archival data interchange.
- Prioritize observability: correlation IDs, structured logs, and actionable metrics.
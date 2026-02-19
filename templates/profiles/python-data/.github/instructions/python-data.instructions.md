---
name: Python Data Standards
description: Typed Python for reliable data-intensive applications
applyTo: "**/*.py"
---

# Python data conventions

- Use type hints for all public functions and methods.
- Use Pydantic models for input/output contracts and schema validation.
- Keep transformations pure and composable where possible.
- Isolate I/O (DB, network, file system) from transformation logic.
- Prefer explicit, testable pipeline steps over monolithic scripts.
- Design for observability: clear logs, metrics, and error context.
- Prefer Parquet for analytical/archival datasets unless requirements conflict.
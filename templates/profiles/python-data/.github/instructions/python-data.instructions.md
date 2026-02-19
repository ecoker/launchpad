---
name: Python Data Systems
description: Typed, composable Python for reliable data-intensive applications
applyTo: "**/*.py"
---

# Python

Python is a tremendous utility. We use it for data pipelines, APIs, scripting,
and ML workflows. But Python's flexibility is also its pitfall — so we constrain
it with types, contracts, and clear structure.

## Type everything

- **Type hints on every public function.** Parameters, return types, class attributes.
  Not optional. Not aspirational. Required.
- **Use Pydantic for data contracts.** Inputs, outputs, configs, API payloads —
  if data crosses a boundary, it goes through a Pydantic model.
- **Run mypy or pyright in CI.** Types that aren't checked are just comments.

```python
# ✅ Typed, validated, explicit
from pydantic import BaseModel
from datetime import datetime

class OrderEvent(BaseModel):
    order_id: str
    customer_id: str
    total_cents: int
    occurred_at: datetime

def process_order(event: OrderEvent) -> dict[str, str]:
    """Process an order event and return a confirmation."""
    return {"order_id": event.order_id, "status": "processed"}


# ❌ Untyped, fragile, mysterious
def process_order(data):
    return {"order_id": data["order_id"], "status": "processed"}
```

## Structure

Organize Python projects with clear separation:

```
src/
  {{PROJECT_NAME}}/
    domain/              # Pure business logic, Pydantic models
      models.py
      pricing.py
    services/            # Orchestration with side effects
      order_service.py
    adapters/            # DB, HTTP clients, queue consumers
      postgres.py
      nats_publisher.py
    api/                 # FastAPI / Flask routes (thin layer)
      routes.py
tests/
  test_pricing.py        # Unit tests against domain logic
  test_order_service.py  # Integration tests
```

## Composable pipelines

- Break transformations into **small, pure functions** that can be composed.
- Each pipeline step should be independently testable with known inputs/outputs.
- Separate I/O from transformation: read data → transform data → write data.

```python
# ✅ Composable pipeline steps
def load_raw_events(path: str) -> list[RawEvent]:
    ...

def validate_events(events: list[RawEvent]) -> list[ValidEvent]:
    return [ValidEvent.model_validate(e) for e in events if e.is_complete()]

def enrich_events(events: list[ValidEvent], lookup: dict[str, Customer]) -> list[EnrichedEvent]:
    return [enrich(e, lookup) for e in events]

def write_parquet(events: list[EnrichedEvent], output_path: str) -> None:
    ...

# Pipeline composition
raw = load_raw_events("s3://bucket/events/")
valid = validate_events(raw)
enriched = enrich_events(valid, customer_lookup)
write_parquet(enriched, "output/enriched.parquet")
```

## Dependencies and tooling

- Use `uv` or `pip-tools` for reproducible dependency management.
- Pin dependencies. Use lockfiles.
- Prefer `pathlib.Path` over string path manipulation.
- Use `structlog` for structured, JSON-formatted logging.
- Use `pytest` with fixtures for testing. Prefer parametrize for edge cases.

## What to avoid

- `**kwargs` as a substitute for proper function signatures.
- Mutable default arguments (`def f(items=[])`).
- Bare `except:` — always catch specific exceptions.
- Global mutable state. If you need shared state, make it explicit.
- `print()` for logging — use a proper logger with structured output.
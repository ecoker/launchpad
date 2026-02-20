---
name: Python + FastAPI
description: Async Python APIs with Pydantic types, ideal for ML/data service backends
applyTo: "**/*.py"
---

# Python + FastAPI

FastAPI is Python's best API framework for AI-assisted development. Pydantic
models serve as both validation and documentation. Async by default.

## Scaffold

```sh
mkdir {{name}} && cd {{name}}
python -m venv .venv
source .venv/bin/activate
pip install fastapi uvicorn[standard] pydantic-settings
```

No official scaffold CLI. Structure the project manually, but follow these
conventions from the start.

## Project structure

```
src/
  main.py                # Entry point — uvicorn target
  config.py              # Pydantic Settings for typed config
  api/
    router.py            # Top-level router aggregation
    dependencies.py      # Shared FastAPI dependencies
    orders/
      router.py          # Route definitions
      schemas.py         # Pydantic request/response models
      service.py         # Business logic
      models.py          # SQLAlchemy/ORM models (if applicable)
    health/
      router.py
  core/
    database.py          # DB session management
    security.py          # Auth utilities
  types/
    common.py            # Shared types
tests/
  conftest.py
  test_orders.py
pyproject.toml           # Project metadata + dependencies
```

## FastAPI patterns

### Pydantic schemas for everything

Every API boundary gets a Pydantic model. This gives you runtime validation,
serialization, and automatic OpenAPI documentation.

```python
# api/orders/schemas.py
from pydantic import BaseModel, Field
from datetime import datetime
from enum import StrEnum


class OrderStatus(StrEnum):
    PENDING = "pending"
    CONFIRMED = "confirmed"
    SHIPPED = "shipped"


class CreateOrderRequest(BaseModel):
    customer_id: str
    items: list["OrderItem"]


class OrderItem(BaseModel):
    product_id: str
    quantity: int = Field(ge=1)


class OrderResponse(BaseModel):
    id: str
    status: OrderStatus
    created_at: datetime

    model_config = {"from_attributes": True}
```

### Route definitions with typed responses

```python
# api/orders/router.py
from fastapi import APIRouter, Depends, HTTPException, status
from .schemas import CreateOrderRequest, OrderResponse
from .service import OrderService

router = APIRouter(prefix="/orders", tags=["orders"])


@router.post("/", response_model=OrderResponse, status_code=status.HTTP_201_CREATED)
async def create_order(
    request: CreateOrderRequest,
    service: OrderService = Depends(),
) -> OrderResponse:
    return await service.create(request)


@router.get("/{order_id}", response_model=OrderResponse)
async def get_order(
    order_id: str,
    service: OrderService = Depends(),
) -> OrderResponse:
    order = await service.get_by_id(order_id)
    if not order:
        raise HTTPException(status_code=404, detail="Order not found")
    return order
```

### Service layer for business logic

```python
# api/orders/service.py
from .schemas import CreateOrderRequest, OrderResponse


class OrderService:
    async def create(self, request: CreateOrderRequest) -> OrderResponse:
        # Business logic here — not in the route handler
        ...

    async def get_by_id(self, order_id: str) -> OrderResponse | None:
        ...
```

### Typed configuration

```python
# config.py
from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    database_url: str
    debug: bool = False
    api_key: str

    model_config = {"env_file": ".env"}


settings = Settings()
```

## Python discipline

- **Type hints everywhere.** Function signatures, return types, variables
  where inference is ambiguous.
- **Pydantic for boundaries.** Every external data shape (API, config, DB)
  gets a Pydantic model.
- **Async by default.** Use `async def` for route handlers. Use sync only
  when interfacing with blocking libraries.
- **Functional core.** Keep business logic in pure functions where possible.
  Push I/O to service boundaries.

## Testing

- **pytest** for everything.
- **httpx + `AsyncClient`** for API integration tests.
- **Unit tests for services.** Test business logic independently from FastAPI.

## What to avoid

- Business logic in route handlers — delegate to services.
- Untyped function signatures.
- `Any` type — use `object` or proper type unions.
- Synchronous blocking calls in async routes.
- Global mutable state — use dependency injection.
- Hand-rolling validation — use Pydantic.

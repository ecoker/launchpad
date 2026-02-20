---
name: Rust + Axum
description: Performance-critical services — Tokio-based, type-safe, zero-cost abstractions
applyTo: "**/*.rs"
---

# Rust + Axum

Rust when correctness and performance are non-negotiable. Axum over Actix —
it's built on Tokio and Tower, composes naturally, and has a simpler
mental model for AI-generated code.

## Scaffold

```sh
cargo new {{name}}
```

Then add dependencies to `Cargo.toml`:

```toml
[dependencies]
axum = "0.7"
tokio = { version = "1", features = ["full"] }
serde = { version = "1", features = ["derive"] }
serde_json = "1"
tower-http = { version = "0.5", features = ["cors", "trace"] }
tracing = "0.1"
tracing-subscriber = "0.3"
```

No official scaffold CLI. Cargo creates the project structure; add
dependencies explicitly.

## Project structure

```
src/
  main.rs                # Entry point — wiring only
  lib.rs                 # Library root
  config.rs              # Typed configuration
  routes/
    mod.rs               # Router aggregation
    orders.rs            # Order routes
    health.rs            # Health check
  services/
    mod.rs
    order_service.rs     # Business logic
  models/
    mod.rs
    order.rs             # Domain types
  db/
    mod.rs
    postgres.rs          # Database implementation
  error.rs               # Error types
  state.rs               # App state
```

## Axum patterns

### Router composition

```rust
// routes/mod.rs
use axum::Router;
use crate::state::AppState;

mod orders;
mod health;

pub fn router() -> Router<AppState> {
    Router::new()
        .merge(health::router())
        .nest("/api/orders", orders::router())
}
```

### Typed extractors and handlers

```rust
// routes/orders.rs
use axum::{
    extract::{Path, State},
    http::StatusCode,
    routing::{get, post},
    Json, Router,
};
use crate::{
    models::order::{CreateOrderRequest, OrderResponse},
    services::order_service::OrderService,
    state::AppState,
    error::AppError,
};

pub fn router() -> Router<AppState> {
    Router::new()
        .route("/", post(create_order))
        .route("/{id}", get(get_order))
}

async fn create_order(
    State(state): State<AppState>,
    Json(request): Json<CreateOrderRequest>,
) -> Result<(StatusCode, Json<OrderResponse>), AppError> {
    let order = state.order_service.create(request).await?;
    Ok((StatusCode::CREATED, Json(order)))
}

async fn get_order(
    State(state): State<AppState>,
    Path(id): Path<String>,
) -> Result<Json<OrderResponse>, AppError> {
    let order = state.order_service.get_by_id(&id).await?;
    Ok(Json(order))
}
```

### Domain types with Serde

```rust
// models/order.rs
use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CreateOrderRequest {
    pub customer_id: String,
    pub items: Vec<OrderItem>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct OrderItem {
    pub product_id: String,
    pub quantity: u32,
}

#[derive(Debug, Clone, Serialize)]
pub struct OrderResponse {
    pub id: String,
    pub status: OrderStatus,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "snake_case")]
pub enum OrderStatus {
    Pending,
    Confirmed,
    Shipped,
}
```

### Error handling

```rust
// error.rs
use axum::{http::StatusCode, response::IntoResponse, Json};
use serde_json::json;

#[derive(Debug)]
pub enum AppError {
    NotFound(String),
    BadRequest(String),
    Internal(anyhow::Error),
}

impl IntoResponse for AppError {
    fn into_response(self) -> axum::response::Response {
        let (status, message) = match self {
            Self::NotFound(msg) => (StatusCode::NOT_FOUND, msg),
            Self::BadRequest(msg) => (StatusCode::BAD_REQUEST, msg),
            Self::Internal(err) => {
                tracing::error!("Internal error: {err:?}");
                (StatusCode::INTERNAL_SERVER_ERROR, "Internal server error".into())
            }
        };
        (status, Json(json!({ "error": message }))).into_response()
    }
}
```

## Rust discipline

- **Explicit error types.** Use `thiserror` or manual error enums. No
  `.unwrap()` in production code.
- **Owned types at boundaries, references internally.** Handlers take
  owned inputs; services can borrow.
- **Keep lifetimes simple.** If lifetime annotations are getting complex,
  clone or restructure.
- **Small traits.** Define traits where consumed, not where implemented.
- **`clippy` as gate.** All clippy warnings are errors in CI.

## Testing

- **Unit tests in `#[cfg(test)]` modules** alongside the code they test.
- **Integration tests in `tests/`** using `axum::test` helpers.
- **Use `tokio::test`** for async tests.

## What to avoid

- `.unwrap()` in production code — use `?` and proper error types.
- Overly complex lifetime annotations — simplify the design instead.
- `unsafe` without documented invariants.
- Giant `main.rs` — extract into modules early.
- Premature optimization — write correct code first, profile second.

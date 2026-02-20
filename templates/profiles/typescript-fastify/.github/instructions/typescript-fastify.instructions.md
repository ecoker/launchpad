---
name: TypeScript + Fastify
description: Schema-driven Node.js APIs with typed routes and plugin architecture
applyTo: "**/*.{ts,js}"
---

# TypeScript + Fastify

Fastify over Express. Always. Fastify enforces structure through its schema
and plugin system, which means AI agents produce more correct code.

## Scaffold

```sh
npm init -y
npm install fastify @fastify/type-provider-typebox @sinclair/typebox
npm install -D typescript @types/node tsx
```

No official scaffold CLI. Initialize manually, but use Fastify's plugin
structure from the start.

## Project structure

```
src/
  server.ts              # Entry point — wiring only
  app.ts                 # Fastify instance + plugin registration
  plugins/
    db.ts                # Database plugin
    auth.ts              # Auth plugin
  routes/
    orders/
      index.ts           # Route plugin registration
      schemas.ts         # TypeBox schemas for this domain
      handlers.ts        # Route handlers (thin)
    health/
      index.ts
  services/
    order-service.ts     # Business logic
  types/
    domain.ts            # Shared domain types
```

## Route patterns

### Schema-first routes

Every route gets a TypeBox schema. This gives you runtime validation AND
TypeScript inference from the same definition.

```typescript
// routes/orders/schemas.ts
import { Type, Static } from '@sinclair/typebox';

export const CreateOrderSchema = {
  body: Type.Object({
    customerId: Type.String({ format: 'uuid' }),
    items: Type.Array(Type.Object({
      productId: Type.String(),
      quantity: Type.Integer({ minimum: 1 }),
    })),
  }),
  response: {
    201: Type.Object({
      id: Type.String({ format: 'uuid' }),
      status: Type.Literal('created'),
    }),
  },
};

export type CreateOrderBody = Static<typeof CreateOrderSchema.body>;
```

### Plugin-based route registration

```typescript
// routes/orders/index.ts
import { FastifyPluginAsync } from 'fastify';
import { CreateOrderSchema } from './schemas';
import { createOrder } from './handlers';

const ordersRoutes: FastifyPluginAsync = async (fastify) => {
  fastify.post('/', { schema: CreateOrderSchema }, createOrder);
};

export default ordersRoutes;
```

### Thin handlers delegating to services

```typescript
// routes/orders/handlers.ts
import { FastifyRequest, FastifyReply } from 'fastify';
import { CreateOrderBody } from './schemas';
import { OrderService } from '../../services/order-service';

export async function createOrder(
  request: FastifyRequest<{ Body: CreateOrderBody }>,
  reply: FastifyReply,
) {
  const order = await OrderService.create(request.body);
  return reply.status(201).send(order);
}
```

## Plugins

Fastify's plugin system is its superpower. Use it for encapsulation:

```typescript
// plugins/db.ts
import fp from 'fastify-plugin';
import { FastifyPluginAsync } from 'fastify';

const dbPlugin: FastifyPluginAsync = async (fastify) => {
  const pool = createPool(fastify.config.DATABASE_URL);
  fastify.decorate('db', pool);

  fastify.addHook('onClose', async () => {
    await pool.end();
  });
};

export default fp(dbPlugin);
```

## TypeScript discipline

- **TypeBox for runtime + compile-time safety.** One schema, two guarantees.
- **Type boundaries explicitly.** Request/response shapes, service interfaces.
- **No `any`.** Use `unknown` and narrow.
- **Discriminated unions** for domain state modeling.

## What to avoid

- Express patterns in Fastify (middleware stacking instead of plugins).
- Business logic in route handlers — keep them thin.
- Skipping schema validation — every route gets a schema.
- `any` types — use TypeBox inference.
- `require()` — use ESM imports.

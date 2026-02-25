---
name: TypeScript + Next.js
description: React ecosystem full-stack with App Router, RSC, and typed boundaries
applyTo: "**/*.{ts,tsx,js,jsx}"
---

# TypeScript + Next.js

Next.js when you specifically need the React ecosystem. App Router with
React Server Components is powerful but has complexity. Be deliberate
about server/client boundaries.

## Scaffold

```sh
npx create-next-app@latest
```

Use the CLI scaffold. Never generate `next.config.js`, `package.json`,
`tsconfig.json`, or other project boilerplate by hand.

## Project structure

```
app/
  layout.tsx             # Root layout (server component)
  page.tsx               # Home page
  (marketing)/           # Route group — no URL segment
    about/page.tsx
  dashboard/
    layout.tsx           # Nested layout
    page.tsx             # Server component — data fetching
    actions.ts           # Server actions
    components/
      chart.tsx          # Client component (use client)
  api/                   # API routes (use sparingly)
lib/
  db.ts                  # Database access
  auth.ts                # Auth utilities
  validations.ts         # Zod schemas
components/
  ui/                    # Shared primitives
types/
  domain.ts              # Shared domain types
```

## Server / Client boundary

This is the hardest part. Get it right:

- **Default to Server Components.** They fetch data, render HTML, stream to
  the client. No JavaScript shipped.
- **`"use client"` only when needed.** Interactive elements: forms, dropdowns,
  modals, charts. Not data display.
- **Keep client components small.** The `"use client"` boundary should wrap
  the smallest interactive surface, not the whole page.
- **Pass server data as props.** Server components fetch, client components
  render interactivity.

```tsx
// ✅ Server component fetches data
// app/dashboard/page.tsx
import { db } from '@/lib/db';
import { DashboardChart } from './components/chart';

export default async function DashboardPage() {
  const metrics = await db.metrics.recent();
  return (
    <div>
      <h1>Dashboard</h1>
      <DashboardChart data={metrics} /> {/* Client component */}
    </div>
  );
}

// ✅ Client component handles interactivity
// app/dashboard/components/chart.tsx
'use client';
import type { Metric } from '@/types/domain';

type Props = { data: Metric[] };

export function DashboardChart({ data }: Props) {
  // Interactive chart rendering
}
```

## Server Actions

- **Use Server Actions for mutations.** They replace API routes for form
  submissions and data mutations.
- **Validate with Zod** at the action boundary.
- **Return structured responses** — not throwing errors for expected failures.

```typescript
// app/dashboard/actions.ts
'use server';

import { z } from 'zod';
import { db } from '@/lib/db';

const CreateProjectSchema = z.object({
  name: z.string().min(1).max(100),
  description: z.string().optional(),
});

export async function createProject(formData: FormData) {
  const parsed = CreateProjectSchema.safeParse(Object.fromEntries(formData));
  if (!parsed.success) {
    return { error: parsed.error.flatten() };
  }
  const project = await db.projects.create(parsed.data);
  return { project };
}
```

## TypeScript discipline

- **`strict: true` in `tsconfig.json`.** Non-negotiable. This enables
  `strictNullChecks`, `noUncheckedIndexedAccess`, `strictFunctionTypes`,
  and all other strict flags. Never disable individual strict checks.
- **Type boundaries explicitly.** API responses, component props, action
  returns — all typed.
- **`type` over `interface`** for domain models — interfaces only for
  contracts implemented by classes (rare in modern TS).
- **Discriminated unions** for state modeling:
  ```typescript
  type AsyncState<T> =
    | { status: "idle" }
    | { status: "loading" }
    | { status: "success"; data: T }
    | { status: "error"; error: string };
  ```
- **No `any`.** Use `unknown` and narrow with type guards.
- **No enums.** Use `as const` objects or string literal union types.
- **Types live next to the code that uses them.** Co-locate types in the
  same file or a sibling `types.ts`. No global `types/` barrel folder.
- **Zod for runtime validation** at server boundaries (Server Actions,
  API routes). Infer TypeScript types from Zod schemas with `z.infer<>`
  to avoid type duplication.
- **Utility types deliberately.** Use `Pick`, `Omit`, `Partial`, `Required`
  to derive types from a single source. Don't duplicate fields across types.

## What to avoid

- Putting interactive code in Server Components without `"use client"`.
- Making everything a Client Component — defeats the purpose of RSC.
- Using API routes for things Server Actions handle better.
- `any` — use `unknown` and narrow.
- Barrel files — they break tree-shaking.
- Client-side data fetching with `useEffect` for initial page data.
- `// @ts-ignore` or `// @ts-expect-error` without a linked issue.

---
name: TypeScript + SvelteKit
description: Full-stack web with intuitive reactivity, SSR, and minimal ceremony
applyTo: "**/*.{ts,js,svelte}"
---

# TypeScript + SvelteKit

SvelteKit is our recommended JS full-stack framework. It has the simplest
mental model of any JS meta-framework — components compile away,
reactivity is built into the language, and the full-stack story is
cohesive. AI agents produce cleaner Svelte code because there's less
framework surface area to get wrong.

## Scaffold

```sh
npm create svelte@latest
```

Use the CLI scaffold. Never generate `package.json`, `svelte.config.js`,
`vite.config.ts`, or other project boilerplate by hand.

## Project structure

Organize by route + feature, not by layer:

```
src/
  routes/
    +layout.svelte       # Root layout
    +page.svelte         # Home page
    dashboard/
      +page.server.ts    # Server load function
      +page.svelte       # Dashboard page
      chart.svelte       # Feature component
  lib/
    server/
      db.ts              # Database (server-only)
      auth.ts            # Auth logic
    components/
      ui/                # Shared presentational components
    stores/
      user.ts            # Domain store
    types/
      domain.ts          # Shared types
  app.html
```

## Svelte patterns

### Components

- **One component, one job.** Small, focused `.svelte` files. If it exceeds
  150 lines, split it.
- **Props are typed.** Use `export let` with TypeScript annotations or
  `$$Props` for complex prop shapes.
- **Events up, data down.** Use `createEventDispatcher` or callback props
  for child-to-parent communication.

```svelte
<script lang="ts">
  export let name: string;
  export let size: 'sm' | 'md' | 'lg' = 'md';

  const sizeClasses = { sm: 'h-8 w-8', md: 'h-10 w-10', lg: 'h-14 w-14' };
</script>

<div class="rounded-full overflow-hidden {sizeClasses[size]}">
  <span>{name[0]}</span>
</div>
```

### Reactivity

- **Keep reactive statements minimal.** `$:` should be simple derivations.
  If the logic is complex, extract it to a function.
- **Avoid deep reactive chains.** If reactive blocks depend on other reactive
  blocks, refactor into explicit computed values.
- **Stores for shared state.** Use Svelte stores (`writable`, `derived`) for
  state shared across components. Keep stores domain-focused.

```typescript
// ✅ Focused, typed store
import { writable, derived } from 'svelte/store';

type User = { id: string; name: string; role: 'admin' | 'user' };

export const currentUser = writable<User | null>(null);
export const isAdmin = derived(currentUser, ($user) => $user?.role === 'admin');
```

### Data loading

- **Use `+page.server.ts` for server-side data.** Load functions run on the
  server and return typed data to the page.
- **Use form actions for mutations.** SvelteKit's `+page.server.ts` actions
  handle POST/PUT/DELETE with progressive enhancement built in.
- **Never fetch in `onMount` for initial data.** Use load functions instead —
  they handle SSR, streaming, and error boundaries automatically.

```typescript
// src/routes/projects/+page.server.ts
import type { PageServerLoad } from './$types';
import { db } from '$lib/server/db';

export const load: PageServerLoad = async () => {
  const projects = await db.projects.findMany();
  return { projects };
};
```

## TypeScript discipline

- **Type boundaries explicitly.** API responses, component props, store
  values — all typed.
- **`type` over `interface`** for domain models.
- **`unknown` over `any`.** Narrow with type guards.
- **Discriminated unions** for state modeling.
- **No enums.** Use `const` objects or union types.

## What to avoid

- `any` — use `unknown` and narrow.
- Client-side fetching for data that should load on the server.
- Global mutable state outside of Svelte stores.
- Barrel files (`index.ts` re-exports) — they break tree-shaking.
- Heavy client-side JS when SvelteKit can handle it server-side.
- Manual WebSocket management — consider whether the use case truly
  needs it or whether server load functions with invalidation suffice.

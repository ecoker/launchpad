---
name: TypeScript + React
description: Strongly typed React applications with composition, clean architecture, and craft
applyTo: "**/*.{ts,tsx,mts,cts,js,jsx}"
---

# TypeScript + React

TypeScript is our strongest opinion. Type everything. Infer where safe.
Never reach for `any`.

## Type discipline

- **Type boundaries explicitly.** API responses, component props, function
  parameters that cross module boundaries — these all get types.
- **Use `type` over `interface`** for domain models. Reserve `interface` for
  things that genuinely need declaration merging (rare).
- **Model state with discriminated unions.** Make illegal states unrepresentable.

```typescript
// ✅ Discriminated union — the compiler enforces correctness
type AsyncState<T> =
  | { status: "idle" }
  | { status: "loading" }
  | { status: "success"; data: T }
  | { status: "error"; error: Error };

// ❌ Bag of optionals — easy to create impossible combinations
type AsyncState<T> = {
  loading: boolean;
  data?: T;
  error?: Error;
};
```

- Prefer `unknown` over `any`. Narrow with type guards.
- Use `as const` and `satisfies` for type-safe config objects.
- Avoid enums — use `const` objects or union types instead.

## React patterns

### Components

- **One component, one job.** If a component has multiple responsibilities,
  extract them. A 200-line component is almost always two or three smaller ones.
- **Props down, events up.** Components receive data via props and communicate
  changes via callbacks. Avoid reaching into children or siblings.
- **Use function components exclusively.** No class components. Ever.

```tsx
// ✅ Small, focused, typed
type UserAvatarProps = {
  name: string;
  imageUrl: string | null;
  size?: "sm" | "md" | "lg";
};

export function UserAvatar({ name, imageUrl, size = "md" }: UserAvatarProps) {
  const sizeClasses = { sm: "h-8 w-8", md: "h-10 w-10", lg: "h-14 w-14" };
  return (
    <div className={`rounded-full overflow-hidden ${sizeClasses[size]}`}>
      {imageUrl ? (
        <img src={imageUrl} alt={name} className="h-full w-full object-cover" />
      ) : (
        <div className="h-full w-full bg-muted flex items-center justify-center text-sm font-medium">
          {name[0]}
        </div>
      )}
    </div>
  );
}
```

### Hooks

- Extract reusable logic into custom hooks. Name them `use<Thing>`.
- Keep hooks focused. A hook that does three things is three hooks.
- Avoid `useEffect` for synchronous derived state — use `useMemo` or just compute it.

### Data loading with React Router v7

- Use **loader functions** for server data. Keep them focused on fetching and
  shaping data for the route.
- Use **action functions** for mutations. They should validate, execute, and redirect.
- Use `useLoaderData()` and `useActionData()` with typed loaders.
- Prefer optimistic UI patterns for mutations that should feel instant.

```tsx
// ✅ Route module with typed loader
export async function loader({ params }: Route.LoaderArgs) {
  const project = await db.projects.findById(params.projectId);
  if (!project) throw new Response("Not found", { status: 404 });
  return { project };
}

export default function ProjectPage({ loaderData }: Route.ComponentProps) {
  const { project } = loaderData;
  return <ProjectDetail project={project} />;
}
```

## Project structure

Organize by feature, not by type:

```
app/
  routes/
    projects/
      route.tsx          # Route module (loader + component)
      project-card.tsx   # Feature-specific component
      use-project.ts     # Feature-specific hook
  components/
    ui/                  # shadcn/ui primitives
    layout/              # Shell, nav, sidebar
  lib/
    db.server.ts         # Database access (server only)
    validation.ts        # Shared schemas (zod, etc.)
  types/
    domain.ts            # Shared domain types
```

## What to avoid

- `any` — use `unknown` and narrow.
- Barrel files (`index.ts` re-exports) — they break tree-shaking and make
  imports harder to trace.
- Default exports for non-route modules — named exports are greppable.
- Prop drilling more than 2 levels — extract context or restructure components.
- `useEffect` for derived/computed values — just compute them inline.
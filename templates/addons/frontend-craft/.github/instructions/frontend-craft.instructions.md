---
name: Frontend Craft
description: Beauty, motion, accessibility, and visual discipline for UI work
applyTo: "**/*.{ts,tsx,js,jsx,css,scss,html,svelte,vue}"
---

# Frontend craft

> "Design is not just what it looks like and feels like.
> Design is how it works." — Steve Jobs

Beautiful software isn't an accident. It's a discipline.

## Visual principles

- **Rhythm and spacing.** Use Tailwind's spacing scale consistently. Every margin,
  padding, and gap should come from the scale — `p-4`, `gap-6`, `mt-8`. Random
  pixel values break visual rhythm.
- **Typography hierarchy.** Establish clear heading levels and body styles.
  Don't reach for `font-bold text-lg` ad hoc — define and reuse type styles.
- **Color with purpose.** Use semantic color tokens: `text-destructive`,
  `bg-muted`, `border-primary`. Tailwind's color palette is a tool, not a buffet.
- **Whitespace is design.** Generous spacing creates clarity. When in doubt,
  add space, don't remove it.

## Tailwind CSS

Tailwind is our styling system. Use it with intention.

```tsx
// ✅ Consistent, readable, responsive
<div className="flex flex-col gap-6 p-6 md:flex-row md:items-center">
  <h2 className="text-2xl font-semibold tracking-tight">Dashboard</h2>
  <p className="text-muted-foreground">Welcome back.</p>
</div>

// ❌ Random values, no rhythm, hard to maintain
<div style={{ padding: '13px', marginTop: '7px', display: 'flex' }}>
  <h2 style={{ fontSize: '19px', fontWeight: 600 }}>Dashboard</h2>
</div>
```

- Prefer utility classes over custom CSS. Extract components, not stylesheets.
- Use `@apply` sparingly and only in component libraries, never in application code.
- Keep responsive design mobile-first: base styles for mobile, `md:` and `lg:` for larger.

## Component design with shadcn/ui

shadcn/ui gives us composable, accessible UI primitives. Use them well.

- **Compose, don't wrap.** Build features by composing `Button`, `Card`,
  `Dialog`, `DropdownMenu` — don't create wrapper components that hide the API.
- **Keep component APIs small.** A component with 15 props needs to be split.
  Prefer compound components and composition.
- **Accessibility is non-negotiable.** shadcn/ui handles most ARIA patterns.
  Don't override them. Add `aria-label` when visual labels aren't sufficient.

```tsx
// ✅ Composable, accessible, clear
<Dialog>
  <DialogTrigger asChild>
    <Button variant="outline">Edit Profile</Button>
  </DialogTrigger>
  <DialogContent>
    <DialogHeader>
      <DialogTitle>Edit your profile</DialogTitle>
      <DialogDescription>Make changes and save when ready.</DialogDescription>
    </DialogHeader>
    <ProfileForm onSubmit={handleSave} />
  </DialogContent>
</Dialog>

// ❌ Monolithic component hiding everything
<EditProfileModal isOpen={open} onClose={close} onSave={save} user={user} />
```

## Motion with purpose

Motion (formerly Framer Motion) brings interfaces to life. Use it to
**clarify**, not to **decorate**.

- **Transitions signal state change.** Fade in new content. Slide out dismissed items.
  Animate height changes so the user's eyes can follow.
- **Keep durations short.** 150–300ms for micro-interactions. 300–500ms for
  layout transitions. Anything longer feels sluggish.
- **Use spring physics** for natural feel: `type: "spring", stiffness: 300, damping: 30`.
- **Respect reduced motion.** Always check `prefers-reduced-motion` and provide
  an instant fallback.

```tsx
// ✅ Purposeful transition that helps comprehension
<motion.div
  initial={{ opacity: 0, y: 8 }}
  animate={{ opacity: 1, y: 0 }}
  exit={{ opacity: 0, y: -8 }}
  transition={{ duration: 0.2 }}
>
  {children}
</motion.div>

// ❌ Motion for its own sake
<motion.div
  animate={{ rotate: 360, scale: [1, 1.5, 1] }}
  transition={{ duration: 2, repeat: Infinity }}
>
  <Button>Click me</Button>
</motion.div>
```

## State management

- **Derive, don't duplicate.** If a value can be computed from other state,
  compute it. Don't store it separately and keep it in sync.
- **Keep state close to where it's used.** Local state > context > global store.
- **URL is state.** Use search params for filters, pagination, and selections
  that should be shareable and bookmarkable.
- **Server state is different from UI state.** Use the right tool for each:
  React Router loaders for server data, local state for UI concerns.

## Performance awareness

- **Measure before optimizing.** Use Lighthouse, Web Vitals, and browser DevTools.
  Don't guess.
- **Lazy-load deliberately.** Split routes. Defer heavy components. Use
  `Suspense` boundaries with meaningful fallbacks.
- **Images matter.** Use modern formats (WebP, AVIF). Set explicit dimensions.
  Lazy-load below-the-fold images.
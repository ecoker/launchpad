---
name: Frontend Craft
description: Visual discipline, interaction quality, and accessibility for every UI surface
applyTo: "**/*.{ts,tsx,js,jsx,css,scss,html,heex,leex,erb,slim,haml,svelte,vue,dart,blade.php}"
---

# Frontend craft

> "Design is not just what it looks like and feels like.
> Design is how it works." — Steve Jobs

Beautiful software isn't an accident. It's a discipline — and it applies
regardless of whether your UI is rendered by React, LiveView, Svelte, Blade,
ERB, or Flutter widgets. These rules are framework-agnostic. The principles are
universal; the implementation adapts to whatever stack you're using.

## Visual principles

- **Rhythm and spacing.** Every margin, padding, and gap should come from a
  consistent spacing scale (4px base: 4, 8, 12, 16, 24, 32, 48, 64). Random
  pixel values break visual rhythm. Use your framework's scale system — Tailwind
  utilities, Flutter `EdgeInsets` constants, or CSS custom properties.
- **Typography hierarchy.** Establish clear heading levels and body styles once,
  then reuse them everywhere. Don't reach for ad-hoc font sizes and weights —
  define a type scale and reference it by name.
- **Color with purpose.** Use semantic color tokens — `primary`, `muted`,
  `destructive`, `success` — not raw hex values scattered across components.
  Define tokens once (CSS variables, Tailwind config, `ThemeData`, or your
  framework's equivalent) and reference them everywhere.
- **Whitespace is design.** Generous spacing creates clarity. When in doubt,
  add space, don't remove it. Cramped layouts feel cheap; breathing room feels
  polished.

## Styling system

**Tailwind CSS is the default styling system for every web framework.** Phoenix,
Rails, SvelteKit, Next.js, Django, and Laravel all have first-class Tailwind
integration. Use it.

- **Utility-first.** Prefer Tailwind utility classes over custom CSS. Extract
  components, not stylesheets. When a pattern repeats, make a component — don't
  reach for `@apply` in a global stylesheet.
- **Theme tokens in `tailwind.config`.** Define your color palette, spacing
  scale, font families, and border radii in the Tailwind config. Every utility
  class then references these tokens automatically. This is how you get
  consistency across hundreds of templates without scattering hex values.
- **CSS custom properties as the bridge.** For server-rendered frameworks
  (Phoenix, Rails, Django, Laravel), define Tailwind's theme tokens as CSS
  custom properties in the root layout (`:root { --color-surface: #1a1a1a; }`).
  This lets non-Tailwind contexts (inline styles, JS, third-party widgets)
  reference the same tokens.
- **Mobile-first responsive design.** Base styles target mobile. Use `md:`,
  `lg:`, `xl:` breakpoints to add complexity for larger screens.
- **`@apply` sparingly.** Only in shared component-level styles where a class
  list would be duplicated across many templates. Never in global CSS.
- **Flutter is the exception.** Flutter doesn't use CSS. Define `ThemeData` and
  `ColorScheme` with palette tokens. Use `ThemeExtension` for custom design
  tokens. Never hardcode colors or spacing — always pull from
  `Theme.of(context)`.

The common thread: **define tokens once, consume by name**.

## Component composition

This is the most important section. Every framework has its own component model
but the discipline is the same.

- **One component, one job.** Whether it's a React component, a LiveView
  function component, a Svelte component, a Rails ViewComponent, or a Flutter
  widget — it should do one thing and do it well.
- **Compose, don't wrap.** Build features by combining small primitives, not by
  creating god-components with 15 props. Prefer compound patterns where the
  consumer controls the structure.
- **Keep components small.** If a component exceeds ~100-150 lines, it's doing
  too much. Extract sub-components.
- **Props down, events up.** Data flows in one direction. Children communicate
  back via callbacks, events, or streams — never by reaching into parent state.
- **Use your framework's component primitives well:**
  - **React / Next.js**: function components, shadcn/ui primitives, compound
    component pattern. Keep client components small; default to server components.
  - **SvelteKit**: `.svelte` components with `$props()`. Shared presentational
    components in `lib/components/ui/`. Use Svelte's built-in reactivity.
  - **Phoenix LiveView**: function components for stateless markup, live
    components for stateful UI. Use `attr` and `slot` for flexible APIs.
    Keep HEEx templates co-located.
  - **Rails**: ViewComponent or Phlex for encapsulated, testable UI primitives.
    Partials for simple shared markup. Stimulus controllers for interaction.
  - **Django**: django-components or template includes. Keep template logic
    minimal — move complexity into template tags or the view layer.
  - **Laravel**: Blade components with clear prop interfaces. Livewire for
    reactive server-rendered UI. Alpine.js for lightweight client interaction.
  - **Flutter**: widget composition over inheritance. Extract widget methods into
    separate widget classes. Use `const` constructors where possible.

## Accessibility

Accessibility is not optional. It's not a Phase 2 feature. It ships with v1.

- **Semantic structure.** Use proper heading hierarchy (`h1` > `h2` > `h3`),
  landmark elements (`nav`, `main`, `aside`), and list elements for lists.
  In Flutter, use `Semantics` widgets for screen reader context.
- **Keyboard navigation.** Every interactive element must be reachable and
  operable via keyboard. Tab order should follow visual order. Custom widgets
  need explicit focus management.
- **ARIA when needed, not by default.** Native HTML elements carry implicit
  roles. Only add `aria-label`, `aria-describedby`, or `role` when the visual
  context isn't accessible to assistive technology.
- **Focus indicators.** Never remove focus outlines without providing an
  alternative. Use visible ring/outline styles on `:focus-visible`.
- **Color is not the only signal.** Use icons, text, or patterns alongside
  color to convey status (success, error, warning).
- **Respect `prefers-reduced-motion`.** Provide instant fallbacks for all
  animations.

## Motion with purpose

Animation clarifies, it doesn't decorate. Use it to help the user understand
what changed and where to look.

- **Transitions signal state change.** Fade in new content. Slide out dismissed
  items. Animate height changes so the user's eyes can follow.
- **Keep durations short.** 150-300ms for micro-interactions. 300-500ms for
  layout transitions. Anything longer feels sluggish.
- **Use your framework's animation system:**
  - **Web (CSS)**: `transition` and `@keyframes` for simple effects. CSS is
    the most performant option and works across all web frameworks.
  - **React**: Motion (formerly Framer Motion) for orchestrated enter/exit.
  - **SvelteKit**: built-in `transition:`, `animate:`, and `in:/out:`
    directives — they're excellent, use them.
  - **Phoenix LiveView**: CSS transitions triggered by `phx-mounted`,
    `phx-remove`, and `phx-*` lifecycle attributes. JS hooks for complex cases.
  - **Rails (Turbo/Stimulus)**: CSS transitions on Turbo frame/stream updates.
    Stimulus controllers for coordinated animation.
  - **Flutter**: `AnimatedContainer`, `Hero`, `AnimationController` for
    explicit control. Keep curves consistent (`Curves.easeInOut` default).
- **Respect `prefers-reduced-motion`.** Always provide an instant fallback.
  This is a hard requirement, not a nice-to-have.

## State management

- **Derive, don't duplicate.** If a value can be computed from other state,
  compute it. Don't store it separately and keep it in sync.
- **Keep state close to where it's used.** Local component state before shared
  stores. The further state travels, the harder it is to reason about.
- **URL is state.** Use search params for filters, pagination, and selections
  that should be shareable and bookmarkable. This applies to every web
  framework — LiveView, SvelteKit, Next.js, Rails, Django, Laravel.
- **Server state is different from UI state.** Don't mix data fetched from the
  server with ephemeral UI concerns (open/closed toggles, hover states, form
  draft text). Use the right tool for each.

## Performance awareness

- **Measure before optimizing.** Use Lighthouse, Web Vitals, browser DevTools,
  or Flutter DevTools. Don't guess where bottlenecks are.
- **Lazy-load deliberately.** Split routes. Defer heavy components. Show
  meaningful loading states, not blank screens.
- **Images matter.** Use modern formats (WebP, AVIF). Set explicit dimensions
  to prevent layout shift. Lazy-load below-the-fold images.
- **Minimize client JavaScript.** Server-rendered frameworks (Phoenix, Rails,
  Django, Laravel) ship minimal JS by design — don't undo that advantage by
  pulling in heavy client bundles. SPA frameworks (React, Svelte) should
  code-split aggressively.

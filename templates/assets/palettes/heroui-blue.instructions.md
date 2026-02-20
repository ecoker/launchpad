# Palette: HeroUI Blue Scale

Use a semantic color system with predictable 50-900 scales and explicit foreground/default values.

## Guidance
- Keep semantic roles: `default`, `primary`, `secondary`, `success`, `warning`, `danger`
- Use lighter defaults in light mode and shifted darker defaults in dark mode
- Keep focus color consistent (`#006fee`) across themes
- Prefer token references over hardcoded hex in component code

## Seed Tokens
- Primary `500`: `#006fee`
- Success `500`: `#17c964`
- Warning `500`: `#f5a524`
- Danger `500`: `#f31260`
- Foreground (light): `#11181c`
- Foreground (dark): `#ecedee`

## Application Rule
For web frameworks, define these tokens in `tailwind.config` under `theme.extend.colors`
using semantic names (`primary`, `success`, `warning`, `danger`). Also set CSS custom
properties on `:root` for non-Tailwind contexts. For Flutter, define them in `ColorScheme`
and `ThemeExtension`. Never scatter raw hex values across components.

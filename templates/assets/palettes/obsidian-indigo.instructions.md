# Palette: Obsidian + Indigo

Dark-first palette inspired by your Phoenix LiveView layout styling.

## Guidance
- Use near-black background and slightly elevated surfaces
- Keep border contrast subtle to avoid visual noise
- Use indigo as the main accent and preserve readable hover states
- Define status colors (`green`, `red`, `yellow`) for realtime state indicators

## Seed Tokens
- Background: `#0f0f0f`
- Surface: `#1a1a1a`
- Surface hover: `#242424`
- Border: `#2a2a2a`
- Text: `#e4e4e7`
- Text muted: `#71717a`
- Accent: `#6366f1`
- Accent hover: `#818cf8`

## Application Rule
For web frameworks, define these tokens in `tailwind.config` under `theme.extend.colors`
and as CSS custom properties on `:root` (e.g. `--color-bg: #0f0f0f; --color-surface: #1a1a1a`).
For Flutter, define in `ThemeData` / `ColorScheme`. Never duplicate literal color values
across components — always reference tokens by name.

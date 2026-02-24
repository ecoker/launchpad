# Palette: Obsidian + Indigo

Dark-first palette with visual energy. Not flat, not boring — polished and alive.

## Guidance
- Use near-black background and slightly elevated surfaces
- Keep border contrast subtle to avoid visual noise
- Use indigo as the main accent and preserve readable hover states
- Define status colors (`green`, `red`, `yellow`) for realtime state indicators
- Apply gradient accents for hero sections and CTAs
- Use soft glows on primary buttons and accent elements

## Seed Tokens
- Background: `#0f0f0f`
- Surface: `#1a1a1a`
- Surface elevated: `#1e1e1e`
- Surface hover: `#242424`
- Border: `#2a2a2a`
- Border hover: `#3a3a3a`
- Text: `#e4e4e7`
- Text muted: `#71717a`
- Accent: `#6366f1`
- Accent hover: `#818cf8`
- Accent glow: `rgba(99, 102, 241, 0.3)`
- Gradient primary: `linear-gradient(135deg, #6366f1, #8b5cf6)`
- Gradient hero: `linear-gradient(180deg, #1a1033 0%, #0f0f0f 100%)`
- Success: `#22c55e`
- Warning: `#eab308`
- Danger: `#ef4444`

## Gradient & Glow Patterns
- Hero sections: use `gradient-hero` as background with accent gradient overlays
- CTA buttons: use `gradient-primary` background with `accent-glow` as box-shadow
- Feature cards: elevated surface with subtle border, gradient accent strip on top or left
- Icon badges: small rounded containers with `accent/10` background tint
- Section transitions: gradient fade between dark and slightly lighter backgrounds

## Application Rule
For web frameworks, define these tokens in `tailwind.config` under `theme.extend.colors`
and as CSS custom properties on `:root` (e.g. `--color-bg: #0f0f0f; --color-surface: #1a1a1a;
--gradient-primary: linear-gradient(135deg, #6366f1, #8b5cf6)`).
For Flutter, define in `ThemeData` / `ColorScheme`. Never duplicate literal color values
across components — always reference tokens by name.

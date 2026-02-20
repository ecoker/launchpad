# Font Pairing: Inter + JetBrains Mono

## Guidance
- Use Inter (or system sans fallback) for UI text and forms
- Use JetBrains Mono (or equivalent monospace fallback) for code, logs, IDs, and diagnostics
- Keep base font-size around 15-16px with line-height ~1.5 for readability
- Avoid introducing more than two primary font families in starter scaffolds

## Application Rule
For web frameworks, define font families in `tailwind.config` under `theme.fontFamily`
(e.g. `sans: ['Inter', ...systemSans]`, `mono: ['JetBrains Mono', ...systemMono]`). Also
set CSS custom properties (`--font-sans`, `--font-mono`) on `:root` for non-Tailwind contexts.
For Flutter, set `fontFamily` in `ThemeData` and reference via `Theme.of(context)`. Never
use raw font-family strings in component code.

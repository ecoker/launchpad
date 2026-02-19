# Font Pairing: Inter + JetBrains Mono

## Guidance
- Use Inter (or system sans fallback) for UI text and forms
- Use JetBrains Mono (or equivalent monospace fallback) for code, logs, IDs, and diagnostics
- Keep base font-size around 15-16px with line-height ~1.5 for readability
- Avoid introducing more than two primary font families in starter scaffolds

## Application Rule
Set shared font tokens early (e.g. CSS variables or theme config) and reference those tokens everywhere.

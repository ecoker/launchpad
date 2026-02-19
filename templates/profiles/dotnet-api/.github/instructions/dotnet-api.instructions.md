---
name: .NET API Standards
description: Clean architecture in C# services with explicit contracts
applyTo: "**/*.{cs,csproj,sln}"
---

# C#/.NET conventions

- Keep application, domain, and infrastructure concerns separated.
- Use explicit DTOs and mapping at service boundaries.
- Avoid anemic models where behavior belongs to the domain.
- Prefer small, focused services over oversized manager classes.
- Use async APIs thoughtfully and propagate cancellation tokens.
- Keep dependency injection composition explicit and test-friendly.
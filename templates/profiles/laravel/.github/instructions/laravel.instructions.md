---
name: Laravel Standards
description: Laravel conventions with clear domain and service boundaries
applyTo: "**/*.{php,blade.php}"
---

# Laravel conventions

- Follow Laravel conventions unless a clear reason exists to deviate.
- Keep controllers thin and delegate business logic to domain/services.
- Use form requests, policies, and validation rules explicitly.
- Avoid fat models with unrelated responsibilities.
- Keep query logic readable and isolated when complex.
- Prefer explicit tests for behavior at HTTP and domain layers.
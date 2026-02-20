---
name: Dart + Flutter
description: Cross-platform native apps — single codebase, widget composition, platform channels
applyTo: "**/*.dart"
---

# Dart + Flutter

Flutter over React Native. Always. Single language, single framework,
no bridge. Widget composition is intuitive for AI agents, and the
rendering engine is consistent across platforms.

## Scaffold

```sh
flutter create {{name}}
```

Use the Flutter CLI. Never hand-write `pubspec.yaml`, platform configs,
or build files.

## Project structure

Organize by feature, not by widget type:

```
lib/
  main.dart                # Entry point — wiring only
  app.dart                 # MaterialApp configuration
  features/
    orders/
      orders_screen.dart   # Screen widget
      order_card.dart      # Feature-specific widget
      order_model.dart     # Domain model
      order_service.dart   # Data access / API
      order_state.dart     # State management
    auth/
      login_screen.dart
      auth_service.dart
      auth_state.dart
  shared/
    widgets/               # Reusable widgets
    theme/
      app_theme.dart       # Theme configuration
      colors.dart          # Color tokens
    services/
      api_client.dart      # HTTP client
      storage.dart         # Local storage
    models/
      user.dart            # Shared domain models
test/
  features/
    orders/
      order_screen_test.dart
```

## Widget patterns

### Small, focused widgets

- **One widget, one job.** If a widget exceeds 100 lines, split it.
- **Prefer composition over inheritance.** Build complex UI by composing
  small widgets, not by extending base classes.
- **Extract widget methods into separate widgets.** Methods that return
  widgets should be separate widget classes for rebuild optimization.

```dart
// ✅ Small, focused, reusable
class OrderCard extends StatelessWidget {
  final Order order;
  final VoidCallback? onTap;

  const OrderCard({super.key, required this.order, this.onTap});

  @override
  Widget build(BuildContext context) {
    return Card(
      child: ListTile(
        title: Text(order.title),
        subtitle: Text(order.status.label),
        trailing: Text(order.total.toStringAsFixed(2)),
        onTap: onTap,
      ),
    );
  }
}
```

### State management

- **Use Riverpod or BLoC** — pick one, commit to it.
- **Keep state classes immutable.** Use `copyWith` or `freezed` for
  state updates.
- **Separate UI state from domain logic.** The widget tree renders state;
  it doesn't compute it.

```dart
// ✅ Immutable state with clear transitions
class OrderState {
  final List<Order> orders;
  final bool isLoading;
  final String? error;

  const OrderState({
    this.orders = const [],
    this.isLoading = false,
    this.error,
  });

  OrderState copyWith({
    List<Order>? orders,
    bool? isLoading,
    String? error,
  }) {
    return OrderState(
      orders: orders ?? this.orders,
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }
}
```

### Data models

- **Immutable data classes.** Use `freezed` or manual `copyWith` patterns.
- **Type everything.** Dart's type system is strong — use it.
- **JSON serialization with `json_serializable`.** Don't hand-write
  `fromJson` / `toJson`.

```dart
// ✅ Typed domain model
class Order {
  final String id;
  final String title;
  final OrderStatus status;
  final double total;
  final DateTime createdAt;

  const Order({
    required this.id,
    required this.title,
    required this.status,
    required this.total,
    required this.createdAt,
  });
}

enum OrderStatus {
  pending('Pending'),
  confirmed('Confirmed'),
  shipped('Shipped');

  final String label;
  const OrderStatus(this.label);
}
```

## Dart discipline

- **Strong typing.** Avoid `dynamic`. Use proper types and generics.
- **Null safety.** Use `?` for nullable types, `!` only when provably non-null.
- **`const` constructors** for immutable widgets — enables rebuild optimization.
- **Named parameters** for constructors with more than 2 arguments.

## Design tokens

- **Theme everything.** Colors, spacing, typography come from the `Theme`.
  Never hardcode values.
- **Use `ThemeExtension`** for custom design tokens beyond Material's defaults.

## What to avoid

- `dynamic` types — always use proper types.
- Massive `build` methods — extract widgets.
- Business logic in widgets — keep it in state/service classes.
- Platform-specific code without abstraction — use platform channels.
- `setState` for complex state — use a proper state management solution.
- Hardcoded colors, fonts, or spacing — use the theme system.

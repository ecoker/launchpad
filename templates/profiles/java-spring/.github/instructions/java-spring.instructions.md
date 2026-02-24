---
name: Java + Spring Boot
description: Clean enterprise Java — DI, auto-configuration, and structured service architecture
applyTo: "**/*.{java,kt,properties,yml,yaml,gradle,xml}"
---

# Java + Spring Boot

Spring Boot when you need enterprise-grade reliability, massive ecosystem
integration, and structured convention. Modern Java (17+) is genuinely
pleasant — records, sealed classes, pattern matching, and text blocks make it
far more expressive than its reputation suggests.

## Scaffold

```sh
spring init --dependencies=web,data-jpa,validation {{name}}
```

Or use `start.spring.io` with: Spring Web, Spring Data JPA, Validation,
Spring Boot Actuator. Gradle (Kotlin DSL) over Maven when possible.

## Project structure

```
src/main/java/com/example/myapp/
  MyAppApplication.java         # Entry point — @SpringBootApplication
  config/
    SecurityConfig.java         # Security configuration
    WebConfig.java              # CORS, interceptors
  order/
    Order.java                  # Domain entity
    OrderService.java           # Business logic
    OrderRepository.java        # Spring Data interface
    OrderController.java        # REST controller (thin)
    OrderDto.java               # Request/response DTOs (records)
  customer/
    Customer.java
    CustomerService.java
    CustomerRepository.java
    CustomerController.java
src/main/resources/
  application.yml               # Configuration
  application-dev.yml           # Dev profile overrides
  application-prod.yml          # Prod profile overrides
src/test/java/com/example/myapp/
  order/
    OrderServiceTest.java
    OrderControllerTest.java
```

- **Organize by domain, not by layer.** `order/` contains entity, service,
  repo, controller, and DTOs together. Don't create `controllers/`, `services/`,
  `repositories/` top-level packages.
- **The application class is wiring only.** No business logic.
- **Configuration classes are explicit.** Use `@Configuration` beans over
  scattering `@Value` annotations through business code.

## Modern Java patterns

Use what the language gives you:

```java
// ✅ Records for immutable DTOs
public record CreateOrderRequest(
    @NotBlank String customerId,
    @NotEmpty List<LineItem> items
) {}

public record OrderSummary(
    String orderId,
    BigDecimal total,
    Instant createdAt
) {}

// ✅ Sealed interfaces for domain modeling
public sealed interface PaymentResult
    permits PaymentResult.Success, PaymentResult.Declined, PaymentResult.Error {
    record Success(String transactionId) implements PaymentResult {}
    record Declined(String reason) implements PaymentResult {}
    record Error(Exception cause) implements PaymentResult {}
}

// ✅ Pattern matching (Java 21+)
public String formatResult(PaymentResult result) {
    return switch (result) {
        case PaymentResult.Success s -> "Paid: " + s.transactionId();
        case PaymentResult.Declined d -> "Declined: " + d.reason();
        case PaymentResult.Error e -> "Error: " + e.cause().getMessage();
    };
}
```

- Use `record` types for all DTOs, API request/response shapes, and value objects.
- Use sealed interfaces for domain types with a fixed set of variants.
- Prefer pattern matching over `instanceof` chains.
- Use text blocks for multi-line strings (SQL, JSON templates).

## Dependency injection

Spring's DI is powerful. Keep it clean:

```java
// ✅ Constructor injection (implicit @Autowired for single constructor)
@Service
public class OrderService {
    private final OrderRepository orders;
    private final PaymentGateway payments;

    public OrderService(OrderRepository orders, PaymentGateway payments) {
        this.orders = orders;
        this.payments = payments;
    }
}

// ❌ Field injection — untestable, hides dependencies
@Service
public class OrderService {
    @Autowired private OrderRepository orders;
}
```

- **Constructor injection always.** Single constructor doesn't need `@Autowired`.
- **Interfaces for boundaries.** Services depend on interfaces, not concrete
  implementations. This is where Spring DI shines.
- **Keep the bean graph simple.** If you need to draw a diagram to understand
  your DI wiring, you have too many beans.

## Controllers

```java
@RestController
@RequestMapping("/api/orders")
public class OrderController {
    private final OrderService orderService;

    public OrderController(OrderService orderService) {
        this.orderService = orderService;
    }

    @PostMapping
    public ResponseEntity<OrderSummary> create(
            @Valid @RequestBody CreateOrderRequest request) {
        var summary = orderService.create(request);
        return ResponseEntity.status(HttpStatus.CREATED).body(summary);
    }

    @GetMapping("/{id}")
    public OrderSummary findById(@PathVariable String id) {
        return orderService.findById(id);
    }
}
```

- **Controllers are thin.** Validate, delegate, shape response. No business logic.
- **Use `@Valid` for request validation.** Jakarta Bean Validation handles it.
- **Return DTOs, not entities.** Never expose JPA entities in API responses.
- **Use `ResponseEntity` when you need status control.** Plain return for 200 OK.

## Data access

```java
// ✅ Spring Data — let Spring generate the implementation
public interface OrderRepository extends JpaRepository<Order, String> {
    List<Order> findByCustomerIdOrderByCreatedAtDesc(String customerId);

    @Query("SELECT o FROM Order o WHERE o.status = :status AND o.createdAt > :since")
    List<Order> findRecentByStatus(@Param("status") OrderStatus status,
                                   @Param("since") Instant since);
}
```

- **Spring Data query methods for simple cases.** Method naming conventions.
- **`@Query` for anything complex.** Don't stretch method names into paragraphs.
- **Entities are not DTOs.** Map between them explicitly in the service layer.
- **Use `application.yml` profiles** for configuration: `application-dev.yml`,
  `application-prod.yml`. Never hardcode database URLs.

## Error handling

```java
@RestControllerAdvice
public class GlobalExceptionHandler {
    @ExceptionHandler(NotFoundException.class)
    public ResponseEntity<ProblemDetail> handleNotFound(NotFoundException ex) {
        var problem = ProblemDetail.forStatusAndDetail(
            HttpStatus.NOT_FOUND, ex.getMessage());
        return ResponseEntity.status(HttpStatus.NOT_FOUND).body(problem);
    }
}
```

- **Use RFC 7807 `ProblemDetail`** (built into Spring 6+) for error responses.
- **Centralize error handling** in `@RestControllerAdvice`.
- **Custom exceptions for domain errors.** `OrderNotFoundException`, not generic
  `RuntimeException`.
- **Never swallow exceptions.** Log meaningful context and propagate.

## Testing

```java
// Unit test — plain JUnit, no Spring context
@Test
void shouldCalculateOrderTotal() {
    var items = List.of(new LineItem("A", 2, new BigDecimal("10.00")));
    var order = Order.create("customer-1", items);
    assertThat(order.total()).isEqualByComparingTo("20.00");
}

// Integration test — slice test for repository
@DataJpaTest
class OrderRepositoryTest {
    @Autowired OrderRepository orders;

    @Test
    void shouldFindByCustomerId() { /* ... */ }
}

// API test — slice test for controller
@WebMvcTest(OrderController.class)
class OrderControllerTest {
    @Autowired MockMvc mockMvc;
    @MockitoBean OrderService orderService;

    @Test
    void shouldCreateOrder() throws Exception { /* ... */ }
}
```

- **Unit tests don't start Spring.** Test business logic with plain objects.
- **Slice tests for infrastructure.** `@DataJpaTest`, `@WebMvcTest` load only
  what's needed. Fast feedback.
- **`@SpringBootTest` sparingly.** Only for full integration tests. They're slow.
- **Use AssertJ** over plain assertions. Fluent, readable, better error messages.

## Configuration

- **`application.yml` over `.properties`.** Structured, readable.
- **Profiles for environments.** `spring.profiles.active=dev` loads
  `application-dev.yml`. Use it.
- **`@ConfigurationProperties`** for typed config binding. Don't scatter
  `@Value("${...}")` across services.
- **Actuator for observability.** Enable health, metrics, info endpoints.
  Don't build custom health checks when Actuator provides them.

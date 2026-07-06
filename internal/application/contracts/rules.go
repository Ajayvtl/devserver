package contracts

// Application Rules & Conventions
//
// 1. Context Propagation:
//    - Every Use Case, Service, and Orchestrator method MUST take `ctx context.Context` as its first parameter.
//    - The Context carries cancellation, deadlines, tracing, and identity.
//
// 2. Transaction Boundaries:
//    - Only Use Cases and Application Services initiate Transactions.
//    - Repositories NEVER commit or rollback; they participate in the `repository.UnitOfWork` injected via context or passed explicitly.
//
// 3. Event Publishing Conventions:
//    - Events are published ONLY after a successful UnitOfWork commit.
//    - Repositories DO NOT publish events. Use `application.EventPublisher`.
//
// 4. Repository Usage:
//    - Repositories are accessed ONLY via their Interfaces defined in `internal/repository`.
//
// 5. Executor Usage:
//    - Executors are requested via the `executor.ExecutorRegistry`. Application logic should never instantiate `local.New()` directly.
//
// 6. Idempotency Rules:
//    - All state-mutating Use Cases must be idempotent where possible. Retries of identical inputs should yield the same domain state without duplicate side-effects.

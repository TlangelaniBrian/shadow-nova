# Database Transactions - Quick Reference

## When to Use Transactions

✅ **USE** transactions for:
- Multiple related database writes
- Operations where partial success is invalid
- Creating hierarchical data (parent + children)
- Updates across multiple tables

❌ **DON'T USE** transactions for:
- Single row operations
- Read-only queries
- Operations with external API calls
- Long-running background jobs

## Basic Pattern

```go
func (s *service) MyOperation(ctx context.Context, data *Model) error {
    return s.WithTx(ctx, func(tx pgx.Tx) error {
        // All DB operations use tx instead of s.db

        err := tx.QueryRow(ctx, query1, ...).Scan(&result1)
        if err != nil {
            return err // Automatic rollback
        }

        _, err = tx.Exec(ctx, query2, ...)
        if err != nil {
            return err // Automatic rollback
        }

        return nil // Automatic commit
    })
}
```

## Common Mistakes

### ❌ Wrong: Using s.db inside transaction
```go
s.WithTx(ctx, func(tx pgx.Tx) error {
    s.db.Exec(ctx, query) // Won't be part of transaction!
})
```

### ✅ Right: Using tx parameter
```go
s.WithTx(ctx, func(tx pgx.Tx) error {
    tx.Exec(ctx, query) // Part of transaction
})
```

### ❌ Wrong: External calls in transaction
```go
s.WithTx(ctx, func(tx pgx.Tx) error {
    tx.Exec(ctx, query)
    sendEmail(user) // Holds lock during I/O!
})
```

### ✅ Right: External calls outside
```go
err := s.WithTx(ctx, func(tx pgx.Tx) error {
    return tx.Exec(ctx, query)
})
if err != nil {
    return err
}
sendEmail(user) // After transaction commits
```

### ❌ Wrong: Using background context
```go
s.WithTx(ctx, func(tx pgx.Tx) error {
    tx.Exec(context.Background(), query) // Ignores timeout!
})
```

### ✅ Right: Pass context through
```go
s.WithTx(ctx, func(tx pgx.Tx) error {
    tx.Exec(ctx, query) // Respects timeout
})
```

## Configuration

```bash
# .env
DB_TX_TIMEOUT=30s  # Default: 30 seconds
```

Adjust based on environment:
- Development: `60s` (more lenient)
- Production: `30s` (standard)
- High-traffic: `5s` (strict)

## Error Handling

```go
// Automatic rollback on error
return s.WithTx(ctx, func(tx pgx.Tx) error {
    if err := doStep1(tx); err != nil {
        return fmt.Errorf("step1 failed: %w", err) // Rolls back
    }

    if err := doStep2(tx); err != nil {
        return fmt.Errorf("step2 failed: %w", err) // Rolls back
    }

    return nil // Commits if we reach here
})
```

## Real Examples

### GitHub Integration
```go
func (s *service) SaveGitHubToken(ctx context.Context, integration *models.GitHubIntegration) error {
    return s.WithTx(ctx, func(tx pgx.Tx) error {
        // Save integration
        err := tx.QueryRow(ctx, insertQuery, ...).Scan(&integration.ID)
        if err != nil {
            return err
        }

        // Update user
        _, err = tx.Exec(ctx, updateQuery, integration.Username, integration.UserID)
        return err
    })
}
```

### Learning Path Creation
```go
func (s *service) CreateLearningPath(ctx context.Context, path *models.LearningPath) error {
    return s.WithTx(ctx, func(tx pgx.Tx) error {
        // Create path
        err := tx.QueryRow(ctx, pathQuery, ...).Scan(&path.CreatedAt)
        if err != nil {
            return err
        }

        // Create modules
        for i := range path.Modules {
            module := &path.Modules[i]
            err := tx.QueryRow(ctx, moduleQuery, ...).Scan(&module.ID)
            if err != nil {
                return err
            }

            // Create lessons for this module
            for j := range module.Lessons {
                lesson := &module.Lessons[j]
                err := tx.QueryRow(ctx, lessonQuery, ...).Scan(&lesson.ID)
                if err != nil {
                    return err
                }
            }
        }

        return nil
    })
}
```

## Testing

```go
func TestTransactionRollback(t *testing.T) {
    db := setupTestDB(t)

    // Execute with error
    err := db.SaveGitHubToken(ctx, invalidData)

    // Verify rollback
    assert.Error(t, err)
    assertNoDataWasCreated(t, db)
}
```

## Troubleshooting

### "Transaction already started"
- Don't nest `WithTx` calls
- Pass `tx` to helper functions instead

### "Transaction timeout"
- Check `DB_TX_TIMEOUT` setting
- Look for slow queries in transaction
- Move non-DB operations outside transaction

### "Deadlock detected"
- Ensure consistent lock order
- Keep transactions short
- Review concurrent access patterns

## See Also

- [TRANSACTIONS.md](./TRANSACTIONS.md) - Detailed documentation
- [TRANSACTION_IMPLEMENTATION.md](./TRANSACTION_IMPLEMENTATION.md) - Implementation summary

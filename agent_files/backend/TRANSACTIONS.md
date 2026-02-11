# Database Transactions in Shadow Nova

This document outlines the transaction management strategy for Shadow Nova's database operations.

## Overview

Database transactions ensure data consistency by grouping multiple operations into an atomic unit. Either all operations succeed and are committed, or if any fails, all changes are rolled back.

## When to Use Transactions

### Always Use Transactions For

1. **Multi-Step Operations with Dependencies**
   - User registration + profile initialization
   - GitHub integration + username update
   - Learning path creation + modules + lessons
   - Project submission + audit logging

2. **Operations That Must Maintain Consistency**
   - Financial transactions
   - Inventory updates
   - State changes that affect multiple tables

3. **Operations Where Partial Success is Invalid**
   - Creating a learning path without modules
   - Saving GitHub tokens without user record update
   - Bulk imports where all-or-nothing is required

### Don't Use Transactions For

1. **Simple Single-Row Operations**
   - Reading a user by email
   - Updating a single field
   - Deleting a single record

2. **Read-Only Operations**
   - Queries and reports
   - Data aggregation
   - Statistics calculation

3. **Long-Running Operations**
   - File uploads
   - External API calls
   - Background processing

## Implementation

### Using WithTx

The preferred method for transaction management is the `WithTx` helper:

```go
func (s *service) SaveGitHubToken(ctx context.Context, integration *models.GitHubIntegration) error {
    return s.WithTx(ctx, func(tx pgx.Tx) error {
        // All database operations use tx instead of s.db

        // Encrypt tokens
        encryptedToken, err := crypto.Encrypt(integration.AccessToken)
        if err != nil {
            return fmt.Errorf("failed to encrypt token: %w", err)
        }

        // Insert integration
        query := `INSERT INTO github_integrations (...) VALUES (...)`
        err = tx.QueryRow(ctx, query, ...).Scan(&integration.ID)
        if err != nil {
            return fmt.Errorf("failed to save integration: %w", err)
        }

        // Update user record
        updateQuery := `UPDATE users SET github_username = $1 WHERE id = $2`
        _, err = tx.Exec(ctx, updateQuery, integration.Username, integration.UserID)
        if err != nil {
            return fmt.Errorf("failed to update user: %w", err)
        }

        // If we return nil, transaction commits
        // If we return error, transaction rolls back
        return nil
    })
}
```

### Using BeginTx (Advanced)

For more complex scenarios where you need fine-grained control:

```go
func (s *service) ComplexOperation(ctx context.Context) error {
    tx, err := s.BeginTx(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx) // Safe to call even after commit

    // Perform operations...
    if err := doStep1(tx); err != nil {
        return err
    }

    if err := doStep2(tx); err != nil {
        return err
    }

    // Explicitly commit
    return tx.Commit(ctx)
}
```

## Transaction Boundaries

### Good Transaction Boundaries

```go
// Good: Tight, focused transaction
func (s *service) CreateUserWithProfile(ctx context.Context, user *models.User) error {
    return s.WithTx(ctx, func(tx pgx.Tx) error {
        // Create user
        query := `INSERT INTO users (...) VALUES (...) RETURNING id`
        err := tx.QueryRow(ctx, query, ...).Scan(&user.ID)
        if err != nil {
            return err
        }

        // Create default profile
        profileQuery := `INSERT INTO user_profiles (user_id, ...) VALUES ($1, ...)`
        _, err = tx.Exec(ctx, profileQuery, user.ID, ...)
        return err
    })
}
```

### Bad Transaction Boundaries

```go
// Bad: Transaction includes non-database operations
func (s *service) BadCreateUser(ctx context.Context, user *models.User) error {
    return s.WithTx(ctx, func(tx pgx.Tx) error {
        // Database operation
        err := tx.QueryRow(ctx, query, ...).Scan(&user.ID)
        if err != nil {
            return err
        }

        // DON'T: External API call inside transaction
        err = sendWelcomeEmail(user.Email) // Locks transaction during network I/O
        if err != nil {
            return err
        }

        // DON'T: Heavy computation inside transaction
        avatar := generateAvatar(user) // Holds transaction lock unnecessarily

        return nil
    })
}

// Good: External operations outside transaction
func (s *service) GoodCreateUser(ctx context.Context, user *models.User) error {
    // Create user in transaction
    err := s.WithTx(ctx, func(tx pgx.Tx) error {
        query := `INSERT INTO users (...) VALUES (...) RETURNING id`
        return tx.QueryRow(ctx, query, ...).Scan(&user.ID)
    })
    if err != nil {
        return err
    }

    // External operations after transaction commits
    go sendWelcomeEmail(user.Email) // Async, doesn't block transaction

    return nil
}
```

## Error Handling

### Automatic Rollback

The `WithTx` helper automatically rolls back on:
- Returned errors
- Panics (after catching and re-raising)
- Context cancellation

```go
return s.WithTx(ctx, func(tx pgx.Tx) error {
    // Any error here triggers automatic rollback
    if err := doOperation(tx); err != nil {
        return err
    }

    // Panic also triggers rollback
    if invalidState {
        panic("unexpected state")
    }

    return nil // Only commits if we return nil
})
```

### Error Context

Always provide context in error messages:

```go
// Good: Clear error context
err := tx.QueryRow(ctx, query, userID).Scan(&result)
if err != nil {
    return fmt.Errorf("failed to create user profile for user %d: %w", userID, err)
}

// Bad: Generic error
if err != nil {
    return err
}
```

## Common Pitfalls

### 1. Nested Transactions

**Problem**: PostgreSQL doesn't support true nested transactions.

```go
// Bad: Attempting to nest transactions
func (s *service) BadNested(ctx context.Context) error {
    return s.WithTx(ctx, func(tx pgx.Tx) error {
        // This will fail - can't start transaction inside transaction
        return s.CreateUser(ctx, user) // CreateUser also uses WithTx
    })
}

// Good: Pass transaction to helpers
func (s *service) GoodNested(ctx context.Context) error {
    return s.WithTx(ctx, func(tx pgx.Tx) error {
        return s.createUserTx(ctx, tx, user) // Use tx-aware helper
    })
}
```

### 2. Long-Running Transactions

**Problem**: Holding locks too long blocks other operations.

```go
// Bad: Long-running transaction
func (s *service) BadLongTx(ctx context.Context) error {
    return s.WithTx(ctx, func(tx pgx.Tx) error {
        // Get records
        rows, _ := tx.Query(ctx, "SELECT * FROM big_table")

        // DON'T: Process for minutes while holding locks
        for rows.Next() {
            processHeavyComputation(row) // Takes 10 seconds per row
        }

        return nil
    })
}

// Good: Minimal transaction duration
func (s *service) GoodShortTx(ctx context.Context) error {
    // Read outside transaction
    rows, _ := s.db.Query(ctx, "SELECT * FROM big_table")

    processed := processAllRows(rows) // Heavy work outside transaction

    // Quick transaction to save results
    return s.WithTx(ctx, func(tx pgx.Tx) error {
        return saveResults(tx, processed)
    })
}
```

### 3. Context Not Passed Through

**Problem**: Using wrong context breaks timeout behavior.

```go
// Bad: Using background context
return s.WithTx(ctx, func(tx pgx.Tx) error {
    // DON'T: Ignores transaction timeout
    _, err := tx.Exec(context.Background(), query)
    return err
})

// Good: Pass context through
return s.WithTx(ctx, func(tx pgx.Tx) error {
    _, err := tx.Exec(ctx, query) // Uses transaction context with timeout
    return err
})
```

### 4. Mixing Transaction and Pool Operations

**Problem**: Some operations use transaction, others don't.

```go
// Bad: Inconsistent usage
return s.WithTx(ctx, func(tx pgx.Tx) error {
    // Uses transaction
    tx.Exec(ctx, "INSERT INTO table1 ...")

    // DON'T: Bypasses transaction
    s.db.Exec(ctx, "INSERT INTO table2 ...") // Not part of transaction!

    return nil
})

// Good: All operations use transaction
return s.WithTx(ctx, func(tx pgx.Tx) error {
    tx.Exec(ctx, "INSERT INTO table1 ...")
    tx.Exec(ctx, "INSERT INTO table2 ...") // Both in same transaction
    return nil
})
```

## Performance Considerations

### Transaction Timeout

Configure `DB_TX_TIMEOUT` based on your needs:

```bash
# Development: More lenient
DB_TX_TIMEOUT=60s

# Production: Strict timeouts
DB_TX_TIMEOUT=30s

# High-traffic APIs: Very short
DB_TX_TIMEOUT=5s
```

### Lock Contention

Minimize time between acquiring and releasing locks:

1. **Prepare data before transaction**
   ```go
   // Prepare outside transaction
   encrypted := crypto.Encrypt(token)
   validated := validateData(input)

   // Quick transaction
   return s.WithTx(ctx, func(tx pgx.Tx) error {
       return tx.Exec(ctx, query, encrypted, validated)
   })
   ```

2. **Use appropriate isolation levels**
   - Default (Read Committed) is usually sufficient
   - Only use higher isolation when necessary

3. **Batch operations**
   ```go
   // Good: Single transaction for batch
   return s.WithTx(ctx, func(tx pgx.Tx) error {
       batch := &pgx.Batch{}
       for _, item := range items {
           batch.Queue(query, item.Field1, item.Field2)
       }
       return tx.SendBatch(ctx, batch).Close()
   })
   ```

## Testing Transactions

### Test Rollback Behavior

```go
func TestTransactionRollback(t *testing.T) {
    // Setup
    db := setupTestDB(t)

    // Execute operation that should rollback
    err := db.WithTx(ctx, func(tx pgx.Tx) error {
        tx.Exec(ctx, "INSERT INTO users (...) VALUES (...)")
        return errors.New("intentional failure")
    })

    // Assert rollback occurred
    assert.Error(t, err)
    assertUserNotExists(t, db)
}
```

### Test Panic Recovery

```go
func TestTransactionPanicRollback(t *testing.T) {
    db := setupTestDB(t)

    assert.Panics(t, func() {
        db.WithTx(ctx, func(tx pgx.Tx) error {
            tx.Exec(ctx, "INSERT INTO users (...) VALUES (...)")
            panic("test panic")
        })
    })

    // Assert rollback occurred
    assertUserNotExists(t, db)
}
```

## Monitoring

### Metrics to Track

1. **Transaction Duration**
   - Alert on transactions exceeding timeout
   - Identify slow operations

2. **Rollback Rate**
   - High rollback rate indicates issues
   - May need better validation

3. **Lock Wait Time**
   - Long waits suggest contention
   - May need to redesign transaction boundaries

4. **Deadlocks**
   - Track deadlock frequency
   - Review transaction order

## Examples from Shadow Nova

### User Registration
```go
// backend/internal/database/users.go
func (s *service) CreateUser(ctx context.Context, user *models.User) error {
    return s.WithTx(ctx, func(tx pgx.Tx) error {
        query := `INSERT INTO users (...) VALUES (...) RETURNING id`
        err := tx.QueryRow(ctx, query, ...).Scan(&user.ID)
        if err != nil {
            return fmt.Errorf("failed to create user: %w", err)
        }

        // Future: Add profile initialization, welcome queue
        return nil
    })
}
```

### GitHub Integration
```go
// backend/internal/database/projects.go
func (s *service) SaveGitHubToken(ctx context.Context, integration *models.GitHubIntegration) error {
    return s.WithTx(ctx, func(tx pgx.Tx) error {
        // Encrypt tokens
        encryptedToken, err := crypto.Encrypt(integration.AccessToken)
        if err != nil {
            return err
        }

        // Save integration
        query := `INSERT INTO github_integrations (...) VALUES (...)`
        err = tx.QueryRow(ctx, query, ...).Scan(&integration.ID)
        if err != nil {
            return err
        }

        // Update user
        _, err = tx.Exec(ctx, `UPDATE users SET github_username = $1 WHERE id = $2`,
            integration.Username, integration.UserID)
        return err
    })
}
```

### Learning Path Seeding
```go
// backend/internal/database/seed.go
func (s *service) seedPathWithModules(ctx context.Context, path *models.LearningPath, seedFunc func(context.Context, pgx.Tx, string) error) error {
    return s.WithTx(ctx, func(tx pgx.Tx) error {
        // Create path
        query := `INSERT INTO learning_paths (...) VALUES (...)`
        err := tx.QueryRow(ctx, query, ...).Scan(&path.CreatedAt)
        if err != nil {
            return err
        }

        // Create modules and lessons
        return seedFunc(ctx, tx, path.ID)
    })
}
```

## Further Reading

- [PostgreSQL Transaction Isolation](https://www.postgresql.org/docs/current/transaction-iso.html)
- [pgx Transaction Documentation](https://pkg.go.dev/github.com/jackc/pgx/v5#Tx)
- [Database Transaction Best Practices](https://www.postgresql.org/docs/current/tutorial-transactions.html)

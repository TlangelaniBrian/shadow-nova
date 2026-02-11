# Crypto Package

This package provides AES-256-GCM encryption for sensitive data stored in the database, specifically GitHub OAuth tokens.

## Features

- **AES-256-GCM Encryption**: Industry-standard authenticated encryption
- **Secure Key Management**: 256-bit keys loaded from environment variables
- **Unique Nonces**: Each encryption operation uses a cryptographically random nonce
- **Authenticated Encryption**: Protects against tampering with built-in authentication

## Setup

### 1. Generate an Encryption Key

Use the provided script to generate a secure 256-bit encryption key:

```bash
./scripts/generate-encryption-key.sh
```

Or manually:

```bash
openssl rand -base64 32
```

### 2. Configure Environment

Add the generated key to your `.env` file:

```env
ENCRYPTION_KEY=<your-generated-key>
```

**IMPORTANT**:
- Never commit the encryption key to version control
- Use different keys for different environments (dev, staging, production)
- Store production keys securely (e.g., AWS Secrets Manager, HashiCorp Vault)

### 3. Initialize in Application

The crypto package is automatically initialized in `main.go`:

```go
if err := crypto.Init(); err != nil {
    log.Fatalf("Failed to initialize encryption: %v", err)
}
```

## Usage

### Encrypting Data

```go
plaintext := "gho_16C7e42F292c6912E7710c838347Ae178B4a"
encrypted, err := crypto.Encrypt(plaintext)
if err != nil {
    return fmt.Errorf("encryption failed: %w", err)
}
// Store encrypted in database
```

### Decrypting Data

```go
// Retrieve encrypted from database
decrypted, err := crypto.Decrypt(encrypted)
if err != nil {
    return fmt.Errorf("decryption failed: %w", err)
}
// Use decrypted token
```

## Migrating Existing Tokens

If you have existing unencrypted tokens in the database, use the migration tool:

```bash
cd backend
go run cmd/migrate-tokens/main.go
```

The migration tool:
- Detects already-encrypted tokens (idempotent)
- Encrypts unencrypted tokens
- Updates records in place
- Provides detailed progress reporting

**WARNING**: Always backup your database before running migrations!

## Security Considerations

### Key Rotation

To rotate encryption keys:

1. Generate a new key
2. Keep the old key temporarily
3. Decrypt data with old key, re-encrypt with new key
4. Update `ENCRYPTION_KEY` environment variable
5. Securely delete old key

### Best Practices

- **Environment Isolation**: Use separate keys for each environment
- **Access Control**: Limit who can access encryption keys
- **Audit Logging**: Monitor encryption/decryption operations
- **Key Backup**: Securely backup keys in multiple locations
- **Regular Rotation**: Rotate keys periodically (e.g., annually)

### What Gets Encrypted

Currently encrypted fields:
- `github_integrations.access_token`
- `github_integrations.refresh_token`

### Encryption Details

- **Algorithm**: AES-256-GCM (Galois/Counter Mode)
- **Key Size**: 256 bits (32 bytes)
- **Nonce Size**: 96 bits (12 bytes) - randomly generated per operation
- **Authentication Tag**: 128 bits (16 bytes) - automatically included
- **Encoding**: Base64 for database storage

### Why AES-GCM?

- **Authenticated Encryption**: Detects tampering
- **Performance**: Hardware-accelerated on most modern CPUs
- **Standard**: NIST-approved, widely used and tested
- **Secure**: No known practical attacks when used correctly

## Testing

Run the crypto package tests:

```bash
cd backend/internal/crypto
go test -v
```

Tests cover:
- Basic encryption/decryption
- Different data sizes and types
- Error conditions
- Nonce uniqueness
- Invalid inputs

## Troubleshooting

### "ENCRYPTION_KEY environment variable not set"

Solution: Set the `ENCRYPTION_KEY` in your `.env` file.

### "ENCRYPTION_KEY must be 32 bytes"

Solution: Generate a new key with `openssl rand -base64 32`.

### "failed to decrypt access token"

Possible causes:
- Wrong encryption key
- Corrupted data in database
- Key was rotated without migrating data

### Migration fails

- Backup database first
- Check database connectivity
- Verify `ENCRYPTION_KEY` is set correctly
- Run with verbose logging to identify specific issues

## Performance

- Encryption: ~1-2 µs per operation (hardware-accelerated)
- Minimal database impact: ~50% size increase for short tokens
- No impact on query performance (only affects read/write)

## Future Enhancements

Potential improvements:
- [ ] Key rotation support
- [ ] Multiple key versions (for gradual rotation)
- [ ] Audit logging for encryption operations
- [ ] Hardware Security Module (HSM) integration
- [ ] Envelope encryption (encrypt keys with master key)

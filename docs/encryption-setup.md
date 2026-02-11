# GitHub OAuth Token Encryption Setup Guide

This guide walks you through setting up encryption for GitHub OAuth tokens in Shadow Nova.

## Overview

Shadow Nova encrypts GitHub OAuth tokens (access and refresh tokens) before storing them in the database using AES-256-GCM encryption.

## Prerequisites

- OpenSSL (for key generation)
- PostgreSQL database
- Go 1.21 or higher

## Step-by-Step Setup

### 1. Generate Encryption Key

Generate a 256-bit encryption key:

```bash
./scripts/generate-encryption-key.sh
```

This will output a base64-encoded key like:
```
Generating 256-bit encryption key...
a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4==

Add this to your .env file as:
ENCRYPTION_KEY=<generated-key>
```

### 2. Add to Environment Variables

Add the key to your `.env` file:

```env
# Encryption key for sensitive data (32 bytes base64 encoded)
# Generate with: openssl rand -base64 32
ENCRYPTION_KEY=a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6q7r8s9t0u1v2w3x4==
```

**Security Notes:**
- Never commit `.env` files to version control
- Use different keys for development, staging, and production
- Store production keys in secure secret management systems

### 3. Verify Configuration

The application will automatically initialize encryption on startup. Check the logs:

```bash
cd backend
go run main.go
```

You should see no errors related to encryption initialization. If you see an error like:
```
Failed to initialize encryption: ENCRYPTION_KEY environment variable not set
```

Then the key is not properly configured in your `.env` file.

### 4. Migrate Existing Tokens (If Applicable)

If you have existing GitHub integrations in your database with unencrypted tokens, run the migration tool:

```bash
cd backend
go run cmd/migrate-tokens/main.go
```

The migration will:
- Detect existing unencrypted tokens
- Encrypt them using the configured key
- Update the database records
- Provide a summary of encrypted records

**Example output:**
```
Connected to database successfully
Found 5 GitHub integrations to process
Successfully encrypted tokens for record 1 (user 123)
Successfully encrypted tokens for record 2 (user 124)
Record 3 (user 125) appears already encrypted, skipping
...

Migration Summary:
Total records: 5
Encrypted: 4
Already encrypted: 1
Failed: 0

Migration completed successfully!
```

### 5. Test the Integration

Test that encryption is working:

```bash
cd backend/internal/crypto
go test -v
```

All tests should pass.

### 6. Verify in Production

After deploying:

1. Check application logs for initialization errors
2. Test GitHub OAuth flow (connect GitHub account)
3. Verify tokens are encrypted in database:

```sql
-- Connect to your database
psql -U user -d shadownova

-- Check a token (should be base64-encoded, ~100-150 chars)
SELECT user_id, LENGTH(access_token), SUBSTRING(access_token, 1, 50)
FROM github_integrations
LIMIT 1;
```

Encrypted tokens will:
- Be longer than original tokens (100+ characters)
- Be base64-encoded (only A-Z, a-z, 0-9, +, /, =)
- Look like: `kN9xE4v5...` instead of `gho_16C7e...`

## Troubleshooting

### Error: "ENCRYPTION_KEY environment variable not set"

**Cause**: Encryption key not configured

**Solution**:
1. Check `.env` file exists in project root
2. Verify `ENCRYPTION_KEY` is set and uncommented
3. Restart the application

### Error: "ENCRYPTION_KEY must be 32 bytes (256 bits)"

**Cause**: Invalid key length

**Solution**:
1. Regenerate key with: `openssl rand -base64 32`
2. Update `.env` file
3. Restart the application

### Error: "failed to decrypt access token"

**Cause**: Wrong encryption key or corrupted data

**Solution**:
1. Verify correct `ENCRYPTION_KEY` is set
2. Check if key was changed after encrypting data
3. If key was lost, you'll need to:
   - Clear `github_integrations` table
   - Have users reconnect their GitHub accounts

### Migration Tool Issues

**Issue**: "Unable to connect to database"

**Solution**:
1. Check `DATABASE_URL` in `.env`
2. Verify PostgreSQL is running
3. Test connection: `psql <DATABASE_URL>`

**Issue**: Migration reports failures

**Solution**:
1. Check logs for specific error messages
2. Verify database permissions
3. Ensure encryption key is correct
4. Consider restoring from backup and re-running

## Security Best Practices

### Key Management

1. **Separate Keys per Environment**
   - Development: Local key in `.env`
   - Staging: Key in CI/CD secrets
   - Production: Key in secure vault (AWS Secrets Manager, etc.)

2. **Access Control**
   - Limit who can view/modify encryption keys
   - Use IAM roles and policies
   - Audit key access

3. **Key Rotation**
   - Rotate keys annually or after security incidents
   - Follow key rotation procedure (see below)

### Backup and Recovery

1. **Backup Encryption Keys**
   - Store in multiple secure locations
   - Use encrypted backup storage
   - Document key recovery procedures

2. **Database Backups**
   - Regular automated backups
   - Test restoration procedures
   - Ensure encrypted data can be restored

### Monitoring

1. **Alert on Encryption Errors**
   - Monitor application logs
   - Set up alerts for decryption failures
   - Track encryption/decryption metrics

2. **Audit Logging**
   - Log GitHub token access
   - Monitor unusual patterns
   - Maintain audit trail

## Key Rotation Procedure

When you need to rotate encryption keys:

1. **Generate New Key**
   ```bash
   openssl rand -base64 32
   ```

2. **Update Environment with Both Keys**
   ```env
   ENCRYPTION_KEY=<new-key>
   ENCRYPTION_KEY_OLD=<old-key>
   ```

3. **Deploy Code to Support Both Keys**
   (This requires code changes to try both keys during decryption)

4. **Run Re-encryption Tool**
   ```bash
   go run cmd/reencrypt-tokens/main.go
   ```

5. **Verify All Tokens Re-encrypted**
   - Check logs
   - Test GitHub integrations
   - Verify database

6. **Remove Old Key**
   - Update environment to remove `ENCRYPTION_KEY_OLD`
   - Deploy changes
   - Securely delete old key

7. **Monitor for Issues**
   - Watch logs for decryption errors
   - Be prepared to rollback if needed

## Database Schema Notes

The existing schema supports encrypted tokens without changes:

```sql
CREATE TABLE github_integrations (
    -- ...
    access_token VARCHAR(255) NOT NULL,
    refresh_token VARCHAR(255),
    -- ...
);
```

- `VARCHAR(255)` is sufficient for encrypted tokens
- Typical encrypted token length: 100-150 characters
- Original token length: 40-60 characters
- Encryption overhead: ~2-2.5x size increase

No schema migration is required for encryption.

## Additional Resources

- [Crypto Package README](../backend/internal/crypto/README.md)
- [AES-GCM Specification](https://nvlpubs.nist.gov/nistpubs/Legacy/SP/nistspecialpublication800-38d.pdf)
- [OWASP Cryptographic Storage Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Cryptographic_Storage_Cheat_Sheet.html)

## Support

If you encounter issues not covered in this guide:

1. Check application logs for detailed error messages
2. Review the troubleshooting section
3. Ensure all prerequisites are met
4. Verify environment configuration
5. Test in development environment first

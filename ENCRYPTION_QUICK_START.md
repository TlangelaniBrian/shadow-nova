# Encryption Quick Start

Quick reference for GitHub OAuth token encryption in Shadow Nova.

## 🔑 Generate Key (First Time Only)

```bash
./scripts/generate-encryption-key.sh
```

Or manually:
```bash
openssl rand -base64 32
```

## ⚙️ Configure Environment

Add to `.env`:
```env
ENCRYPTION_KEY=<your-generated-key-here>
```

## 🚀 Start Application

```bash
cd backend
go run main.go
```

Should see no errors about encryption.

## 🔄 Migrate Existing Tokens (If Needed)

```bash
cd backend
go run cmd/migrate-tokens/main.go
```

## ✅ Test

```bash
cd backend/internal/crypto
go test -v
```

## 🔍 Verify Tokens Are Encrypted

Connect to database and check:
```sql
SELECT user_id, LENGTH(access_token) as token_length,
       SUBSTRING(access_token, 1, 20) as token_start
FROM github_integrations
LIMIT 1;
```

Encrypted tokens:
- Length: 100-150+ characters
- Start: Base64 characters only (A-Z, a-z, 0-9, +, /, =)
- Example: `kN9xE4v5fG7h8I9...`

Unencrypted tokens:
- Length: 40-60 characters
- Start: `gho_`, `ghr_`, etc.
- Example: `gho_16C7e42F292c...`

## ⚠️ Troubleshooting

| Error | Solution |
|-------|----------|
| "ENCRYPTION_KEY environment variable not set" | Add key to `.env` file |
| "ENCRYPTION_KEY must be 32 bytes" | Regenerate key with correct command |
| "failed to decrypt access token" | Check if key changed, may need to re-run migration |

## 🔒 Security Reminders

- ✅ Never commit `.env` to git
- ✅ Use different keys per environment
- ✅ Store production keys in secure vault
- ✅ Backup keys in multiple secure locations
- ✅ Rotate keys annually

## 📚 More Info

- Full setup: `docs/encryption-setup.md`
- Package docs: `backend/internal/crypto/README.md`
- Implementation: `ENCRYPTION_IMPLEMENTATION.md`

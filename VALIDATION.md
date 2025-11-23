# Shadow Nova - Validation Guide

## Overview

Shadow Nova uses **go-playground/validator** for request validation, providing:

- Declarative validation via struct tags
- Custom validators
- Auto-generated error messages
- Type-safe validation

## Usage Examples

### 1. Basic Request Validation

```go
type RegisterRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Username string `json:"username" validate:"required,min=3,max=20,alphanum"`
    Password string `json:"password" validate:"required,min=8"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req models.RegisterRequest

    if err := validator.ValidateRequest(r, &req); err != nil {
        validator.WriteValidationError(w, err)
        return
    }

    // Process valid request...
}
```

### 2. Available Validation Tags

```go
validate:"required"              // Field must not be empty
validate:"email"                 // Must be valid email
validate:"min=3"                 // Minimum length/value
validate:"max=20"                // Maximum length/value
validate:"alphanum"              // Only letters and numbers
validate:"url"                   // Valid URL
validate:"gte=18"                // Greater than or equal
validate:"lte=100"               // Less than or equal
validate:"oneof=admin user"      // Must be one of values
validate:"omitempty,url"         // Optional but if present must be URL
```

### 3. Custom Validators

We've added a `strong_password` validator:

```go
type ChangePasswordRequest struct {
    Password string `json:"password" validate:"required,min=8,strong_password"`
}
```

This requires:

- Uppercase letter
- Lowercase letter
- Number
- Special character (!@#$%^&\*()\_+-=[]{}|;:,.<>?)

### 4. Testing Validation

**Valid Request:**

```bash
curl -X POST http://localhost:3000/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "johndoe",
    "password": "SecurePass123!"
  }'
```

**Invalid Request (will fail validation):**

```bash
curl -X POST http://localhost:3000/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "not-an-email",
    "username": "ab",
    "password": "weak"
  }'
```

**Response:**

```json
{
  "error": "Validation failed",
  "details": {
    "email": "email must be a valid email address",
    "username": "username must be at least 3 characters",
    "password": "password must be at least 8 characters"
  }
}
```

## Adding Custom Validators

Edit `backend/internal/validator/validator.go`:

```go
func init() {
    validate = validator.New()

    // Register your custom validator
    validate.RegisterValidation("custom_rule", validateCustomRule)
}

func validateCustomRule(fl validator.FieldLevel) bool {
    value := fl.Field().String()
    // Your validation logic
    return true
}
```

Then update `getErrorMessage()` to provide user-friendly messages:

```go
case "custom_rule":
    return field + " must satisfy custom rule"
```

## API Endpoints

### Public Endpoints (No Auth)

- `POST /api/register` - User registration
- `POST /api/login` - User login

### Protected Endpoints (JWT Required)

- `POST /api/progress` - Update learning progress
- `GET /api/paths` - Get learning paths
- `GET /api/projects` - Get projects

## Error Response Format

All validation errors follow this format:

```json
{
  "error": "Validation failed",
  "details": {
    "field_name": "user-friendly error message"
  }
}
```

Success responses:

```json
{
  "message": "Operation successful",
  "data": { ... }
}
```

## Best Practices

1. **Always validate at the API boundary** - Never trust client input
2. **Use specific validators** - `email`, `url`, `alphanum` are better than just `required`
3. **Provide clear error messages** - Update `getErrorMessage()` for custom validators
4. **Keep validation logic in tags** - Don't scatter validation across handlers
5. **Test edge cases** - Empty strings, special characters, max lengths

## Next Steps

- [ ] Add database-level uniqueness checks (email, username)
- [ ] Implement password hashing before storage
- [ ] Add JWT token generation and validation
- [ ] Create integration tests for validation
- [ ] Add rate limiting per endpoint

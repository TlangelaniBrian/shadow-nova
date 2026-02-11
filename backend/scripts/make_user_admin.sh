#!/bin/bash

# Script to make a user an admin
# Usage: ./make_user_admin.sh <email>

if [ -z "$1" ]; then
    echo "Usage: $0 <email>"
    echo "Example: $0 user@example.com"
    exit 1
fi

EMAIL="$1"

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Update user role to admin
psql "$DATABASE_URL" <<EOF
UPDATE users
SET user_role = 'admin'
WHERE email = '$EMAIL';

SELECT id, email, username, user_role
FROM users
WHERE email = '$EMAIL';
EOF

echo ""
echo "User $EMAIL has been updated to admin role (if they exist)"

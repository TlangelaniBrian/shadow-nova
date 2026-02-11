#!/bin/bash
echo "Generating 256-bit encryption key..."
openssl rand -base64 32
echo ""
echo "Add this to your .env file as:"
echo "ENCRYPTION_KEY=<generated-key>"

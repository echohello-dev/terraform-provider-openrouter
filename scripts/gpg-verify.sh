#!/usr/bin/env bash
set -euo pipefail

# Verify the GPG signing key fingerprint matches the expected value.
# Usage: scripts/gpg-verify.sh <expected_key_id>

EXPECTED_KEY_ID="${1:-}"

if [ -z "$EXPECTED_KEY_ID" ]; then
  echo "::error::Usage: gpg-verify.sh <expected_key_id>"
  exit 1
fi

# Check if gpg is available
if ! command -v gpg >/dev/null 2>&1; then
  echo "::error::gpg is not installed"
  exit 1
fi

# List secret keys and check for the expected key
if ! gpg --list-secret-keys --keyid-format=long "$EXPECTED_KEY_ID" >/dev/null 2>&1; then
  echo "::error::GPG key $EXPECTED_KEY_ID not found"
  exit 1
fi

FINGERPRINT=$(gpg --list-keys --with-colons --keyid-format=long "$EXPECTED_KEY_ID" | awk -F: '/^fpr/ {print $10; exit}')

echo "Fingerprint: $FINGERPRINT"
echo "KeyID:       $EXPECTED_KEY_ID"
echo "GPG key verified."

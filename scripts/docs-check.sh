#!/usr/bin/env bash
set -euo pipefail

# Check if documentation is up to date.
# Run this after `mise run doc` or `tfplugindocs generate`.

if [ -n "$(git status --porcelain docs/)" ]; then
  echo "::error::Documentation is not up to date. Run 'mise run doc' and commit the changes."
  echo ""
  echo "Changed files:"
  git status --short docs/
  echo ""
  echo "Diff:"
  git diff docs/
  exit 1
fi

echo "Documentation is up to date."

#!/bin/sh

branch="$(git rev-parse --abbrev-ref HEAD)"

if [ "$branch" = "main" ]; then
  echo "No committing to main!"
  exit 1
fi

# https://prettier.io/docs/en/precommit.html#option-6-shell-script
FILES=$(git diff --cached --name-only --diff-filter=ACMR | sed 's| |\\ |g')
[ -z "$FILES" ] && exit 0

# Prettify all selected files
echo "$FILES" | xargs ./node_modules/.bin/prettier --ignore-unknown --write

# Terraform fmt everything
echo "$FILES" | grep '.tf' | xargs ./bin/terraform fmt

# Add back the modified/prettified files to staging
echo "$FILES" | xargs git add

exit 0

#!/bin/bash

# Run the break-check binary and capture output
output=$(./break-check 2>&1)
exit_code=$?

# Post a comment on the PR
comment_body="Break-Check completed with exit code $exit_code\n\`\`\`\n$output\n\`\`\`"
curl -s -H "Authorization: token $GITHUB_TOKEN" \
     -H "Accept: application/vnd.github.v3+json" \
     -d "{\"body\": \"$comment_body\"}" \
     "https://api.github.com/repos/$GITHUB_REPOSITORY/issues/$PR_NUMBER/comments"

# Exit with 0 to always pass the action
exit 0
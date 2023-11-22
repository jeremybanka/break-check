#!/bin/bash

# Run the break-check binary and capture output
output=$(./break-check 2>&1)
exit_code=$?

# Escape special characters in the output
escaped_output=$(echo "$output" | sed 's/"/\\"/g' | sed ':a;N;$!ba;s/\n/\\n/g')

# Prepare comment body
comment_body="Break-Check completed with exit code $exit_code\n\`\`\`\n$escaped_output\n\`\`\`"

# Debug: Print the comment body
echo "Comment Body:"
echo "$comment_body"

# Post a comment on the PR
curl -s -H "Authorization: token $GITHUB_TOKEN" \
     -H "Accept: application/vnd.github.v3+json" \
     -d "{\"body\": \"$comment_body\"}" \
     "https://api.github.com/repos/$GITHUB_REPOSITORY/issues/$PR_NUMBER/comments"

# Exit with 0 to always pass the action
exit 0

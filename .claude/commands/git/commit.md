---
allowed-tools: Bash(yarn format), Bash(make generate-types), Bash(git add:*), Bash(git commit:*), Bash(git diff:*), Bash(git status:*)
description: Commit all changes
---

# commit

Create a git commit by analyzing all changes and generating an appropriate commit message.

## Usage

```
/git:commit
```

## Implementation

1. Generate types and format code:
   - Run `cd backend && make generate-types` to update OpenAPI types
   - Run `cd ../frontend/svetu && yarn lint:fix && yarn format` to fix linting and formatting
2. Run `git status` to see all untracked and modified files
3. Run `git diff` to see unstaged changes
4. Run `git diff --cached` to see staged changes
5. Add all changes with `git add -A`
6. Analyze the complete diff to understand what changes were made
7. Generate a concise commit message based on the analysis that:
   - Focuses on the "why" rather than the "what"
   - Is 1-2 sentences maximum
   - Does NOT include any Claude attribution or emojis
   - Does NOT mention "Generated with Claude Code" or "Co-Authored-By"
8. Create the commit with the generated message
9. Show the result with `git log -1 --oneline`

## Important

- NEVER add Claude attribution to commits
- Keep commit messages professional and concise
- Focus on describing the actual changes made

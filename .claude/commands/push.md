---
allowed-tools: Bash(yarn format), Bash(git add:*), Bash(git commit:*), Bash(git push:*)
description: Format code, commit all changes, and push to remote
---

## Usage

```
/git:push
```

## Context

Gather context about current state:

- Current git status: !`git status --porcelain`
- Current branch: !`git branch --show-current`
- Modified files: !`git diff --name-only`
- Staged files: !`git diff --cached --name-only`

## Your task

Perform the following sequence:

1. **Format code**: Run `yarn format` to ensure code is properly formatted
2. **Stage changes**: Add all changes to git staging area with `git add .`
3. **Create commit**: Write a concise commit message (1-2 lines max) that describes the changes without Claude authorship
4. **Push changes**: Push to the remote repository

Follow conventional commit format when possible (feat:, fix:, docs:, etc.) but keep it brief and descriptive.

Do not include any authorship attribution to Claude in the commit message.
